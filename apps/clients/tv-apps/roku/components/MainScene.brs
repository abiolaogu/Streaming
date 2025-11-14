sub init()
    print "MainScene: init"

    m.top.backgroundURI = ""
    m.top.backgroundColor = "0x000000"

    ' Get UI components
    m.navigationMenu = m.top.findNode("navigationMenu")
    m.featuredPoster = m.top.findNode("featuredPoster")
    m.contentRowList = m.top.findNode("contentRowList")
    m.videoPlayer = m.top.findNode("videoPlayer")
    m.loginScreen = m.top.findNode("loginScreen")
    m.detailsScreen = m.top.findNode("detailsScreen")
    m.loadingSpinner = m.top.findNode("loadingSpinner")

    ' Set up event handlers
    m.contentRowList.observeField("itemFocused", "onContentItemFocused")
    m.contentRowList.observeField("itemSelected", "onContentItemSelected")
    m.videoPlayer.observeField("state", "onVideoPlayerStateChange")

    ' Initialize app state
    m.currentScreen = "home"
    m.isLoggedIn = false

    ' Check authentication
    checkAuthentication()

    ' Load initial content
    if m.isLoggedIn then
        loadHomeContent()
    else
        showLoginScreen()
    end if
end sub

sub checkAuthentication()
    ' Check if user is authenticated
    sec = CreateObject("roRegistrySection", "StreamVerseAuth")
    if sec.Exists("authToken") then
        m.authToken = sec.Read("authToken")
        m.isLoggedIn = (m.authToken <> "")
    end if
end sub

sub showLoginScreen()
    m.loginScreen.visible = true
    m.contentRowList.visible = false

    ' Set up login button handlers
    loginButton = m.loginScreen.findNode("loginButton")
    loginButton.observeField("buttonSelected", "onLoginButtonPressed")

    signupButton = m.loginScreen.findNode("signupButton")
    signupButton.observeField("buttonSelected", "onSignupButtonPressed")
end sub

sub onLoginButtonPressed()
    emailInput = m.loginScreen.findNode("emailInput")
    passwordInput = m.loginScreen.findNode("passwordInput")

    email = emailInput.text
    password = passwordInput.text

    ' Call authentication API
    performLogin(email, password)
end sub

sub performLogin(email as String, password as String)
    m.loadingSpinner.visible = true

    ' API call to authentication service
    port = CreateObject("roMessagePort")
    urlTransfer = CreateObject("roUrlTransfer")
    urlTransfer.SetPort(port)
    urlTransfer.SetUrl("https://api.streamverse.io/v1/auth/login")
    urlTransfer.SetRequest("POST")
    urlTransfer.AddHeader("Content-Type", "application/json")

    requestBody = {
        email: email,
        password: password
    }
    urlTransfer.SetMessagePort(port)

    if urlTransfer.AsyncPostFromString(FormatJson(requestBody)) then
        while true
            msg = wait(5000, port)
            if type(msg) = "roUrlEvent" then
                if msg.GetResponseCode() = 200 then
                    response = ParseJson(msg.GetString())
                    ' Save auth token
                    sec = CreateObject("roRegistrySection", "StreamVerseAuth")
                    sec.Write("authToken", response.token)
                    sec.Flush()

                    m.authToken = response.token
                    m.isLoggedIn = true

                    ' Hide login screen and load content
                    m.loginScreen.visible = false
                    loadHomeContent()
                else
                    ' Show error
                    print "Login failed: "; msg.GetResponseCode()
                end if
                exit while
            end if
        end while
    end if

    m.loadingSpinner.visible = false
end sub

sub loadHomeContent()
    m.loadingSpinner.visible = true

    ' Create API request
    port = CreateObject("roMessagePort")
    urlTransfer = CreateObject("roUrlTransfer")
    urlTransfer.SetPort(port)
    urlTransfer.SetUrl("https://api.streamverse.io/v1/content/home")
    urlTransfer.AddHeader("Authorization", "Bearer " + m.authToken)
    urlTransfer.SetMessagePort(port)

    if urlTransfer.AsyncGetToString() then
        while true
            msg = wait(5000, port)
            if type(msg) = "roUrlEvent" then
                if msg.GetResponseCode() = 200 then
                    response = ParseJson(msg.GetString())
                    displayHomeContent(response)
                else
                    print "Failed to load content: "; msg.GetResponseCode()
                end if
                exit while
            end if
        end while
    end if

    m.loadingSpinner.visible = false
end sub

sub displayHomeContent(content as Object)
    ' Set featured content
    if content.featured <> invalid then
        m.featuredPoster.uri = content.featured.image
    end if

    ' Create content rows
    contentNode = CreateObject("roSGNode", "ContentNode")

    for each row in content.rows
        rowNode = contentNode.createChild("ContentNode")
        rowNode.title = row.title

        for each item in row.items
            itemNode = rowNode.createChild("ContentNode")
            itemNode.title = item.title
            itemNode.description = item.description
            itemNode.hdPosterUrl = item.thumbnail
            itemNode.streamUrl = item.streamUrl
            itemNode.contentId = item.id
        end for
    end for

    m.contentRowList.content = contentNode
    m.contentRowList.visible = true
    m.contentRowList.setFocus(true)
end sub

sub onContentItemFocused()
    ' Handle item focus for preview
    itemFocused = m.contentRowList.itemFocused
    print "Item focused: "; itemFocused
end sub

sub onContentItemSelected()
    ' Handle item selection
    item = m.contentRowList.content.getChild(m.contentRowList.rowItemFocused[0]).getChild(m.contentRowList.rowItemFocused[1])
    showDetailsScreen(item)
end sub

sub showDetailsScreen(item as Object)
    m.detailsScreen.visible = true
    m.contentRowList.visible = false

    ' Populate details
    m.detailsScreen.findNode("detailsTitle").text = item.title
    m.detailsScreen.findNode("detailsDescription").text = item.description
    m.detailsScreen.findNode("detailsPoster").uri = item.hdPosterUrl

    ' Set up button handlers
    playButton = m.detailsScreen.findNode("playButton")
    playButton.observeField("buttonSelected", "onPlayButtonPressed")
    playButton.setFocus(true)
end sub

sub onPlayButtonPressed()
    ' Get selected content
    item = m.contentRowList.content.getChild(m.contentRowList.rowItemFocused[0]).getChild(m.contentRowList.rowItemFocused[1])
    playVideo(item)
end sub

sub playVideo(item as Object)
    ' Hide other screens
    m.detailsScreen.visible = false
    m.contentRowList.visible = false
    m.videoPlayer.visible = true

    ' Configure video player
    videoContent = CreateObject("roSGNode", "ContentNode")
    videoContent.url = item.streamUrl
    videoContent.title = item.title
    videoContent.streamFormat = "hls"

    ' DRM configuration if needed
    if item.drmLicenseUrl <> invalid then
        videoContent.encodingType = "PlayReadyLicenseAcquisitionUrl"
        videoContent.encodingKey = item.drmLicenseUrl
    end if

    m.videoPlayer.content = videoContent
    m.videoPlayer.control = "play"
    m.videoPlayer.setFocus(true)
end sub

sub onVideoPlayerStateChange()
    state = m.videoPlayer.state
    print "Video player state: "; state

    if state = "finished" or state = "error" then
        ' Return to home
        m.videoPlayer.visible = false
        m.videoPlayer.control = "stop"
        m.contentRowList.visible = true
        m.contentRowList.setFocus(true)
    end if
end sub

function onKeyEvent(key as String, press as Boolean) as Boolean
    handled = false

    if press then
        if key = "back" then
            if m.videoPlayer.visible then
                m.videoPlayer.visible = false
                m.videoPlayer.control = "stop"
                m.contentRowList.visible = true
                m.contentRowList.setFocus(true)
                handled = true
            else if m.detailsScreen.visible then
                m.detailsScreen.visible = false
                m.contentRowList.visible = true
                m.contentRowList.setFocus(true)
                handled = true
            end if
        end if
    end if

    return handled
end function
