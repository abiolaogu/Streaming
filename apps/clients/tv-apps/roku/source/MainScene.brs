sub init()
    m.top.backgroundUri = "pkg:/images/background.jpg"
    m.title = m.top.findNode("title")
    m.contentGrid = m.top.findNode("contentGrid")
    m.loadingSpinner = m.top.findNode("loadingSpinner")
    m.errorMessage = m.top.findNode("errorMessage")
    
    ' Set up focus
    m.contentGrid.setFocus(true)
    
    ' Load content
    loadHomeContent()
end sub

sub loadHomeContent()
    m.loadingSpinner.visible = true
    
    ' API call
    url = "https://api.streamverse.com/api/v1/content/home"
    request = CreateObject("roUrlTransfer")
    request.SetUrl(url)
    request.SetCertificatesFile("common:/certs/ca-bundle.crt")
    request.AddHeader("Authorization", "Bearer " + GetAuthToken())
    request.AddHeader("Accept", "application/json")
    
    response = request.GetToString()
    m.loadingSpinner.visible = false
    
    if response <> invalid
        json = ParseJson(response)
        if json <> invalid
            populateContentGrid(json)
        else
            showError("Failed to parse response")
        end if
    else
        showError("Failed to load content")
    end if
end sub

sub populateContentGrid(contentRows)
    rootContentNode = CreateObject("roSGNode", "ContentNode")
    
    for each row in contentRows
        rowNode = CreateObject("roSGNode", "ContentNode")
        rowNode.title = row.title
        
        for each item in row.items
            contentNode = CreateObject("roSGNode", "ContentNode")
            contentNode.title = item.title
            contentNode.hdPosterUrl = item.posterUrl
            contentNode.url = item.streamUrl
            contentNode.id = item.id
            rowNode.appendChild(contentNode)
        end for
        
        rootContentNode.appendChild(rowNode)
    end for
    
    m.contentGrid.content = rootContentNode
    
    ' Set up selection handler
    m.contentGrid.observeField("itemSelected", "onContentSelected")
end sub

sub onContentSelected(event)
    selectedIndex = event.getData()
    selectedRow = m.contentGrid.rowItemSelected[0]
    selectedItem = m.contentGrid.content.getChild(selectedRow).getChild(selectedIndex[1])
    
    ' Show detail screen or play video
    showContentDetail(selectedItem)
end sub

sub showContentDetail(content)
    ' Navigate to detail screen
    detailScreen = CreateObject("roSGNode", "ContentDetailScreen")
    detailScreen.content = content
    m.top.appendChild(detailScreen)
end sub

sub showError(message)
    m.errorMessage.text = message
    m.errorMessage.visible = true
end sub

function GetAuthToken() as String
    ' Get token from registry
    sec = CreateObject("roRegistrySection", "auth")
    if sec.Exists("token")
        return sec.Read("token")
    end if
    return ""
end function

