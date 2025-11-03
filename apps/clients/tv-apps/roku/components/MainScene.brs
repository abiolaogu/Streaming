'
' Main Scene Controller
'

function Init()
    m.loginView = m.top.findNode("loginView")
    m.homeView = m.top.findNode("homeView")
    m.detailView = m.top.findNode("detailView")
    m.playerView = m.top.findNode("playerView")
    m.searchView = m.top.findNode("searchView")
    
    ' Check authentication
    if GetAuthToken() <> invalid
        ShowHome()
    else
        ShowLogin()
    end if
    
    ' Set up observers
    m.loginView.ObserveField("loginSuccess", "OnLoginSuccess")
    m.homeView.ObserveField("contentSelected", "OnContentSelected")
    m.detailView.ObserveField("playRequested", "OnPlayRequested")
    m.homeView.ObserveField("searchRequested", "OnSearchRequested")
end function

function ShowLogin()
    m.loginView.visible = true
    m.homeView.visible = false
    m.loginView.SetFocus(true)
end function

function ShowHome()
    m.loginView.visible = false
    m.homeView.visible = true
    m.homeView.SetFocus(true)
    m.homeView.LoadContent()
end function

function OnLoginSuccess()
    ShowHome()
end function

function OnContentSelected(msg)
    content = msg.GetData()
    m.homeView.visible = false
    m.detailView.visible = true
    m.detailView.content = content
    m.detailView.SetFocus(true)
end function

function OnPlayRequested(msg)
    content = msg.GetData()
    m.detailView.visible = false
    m.playerView.visible = true
    m.playerView.content = content
    m.playerView.SetFocus(true)
end function

function OnSearchRequested()
    m.homeView.visible = false
    m.searchView.visible = true
    m.searchView.SetFocus(true)
end function

function GetAuthToken() as String
    sec = CreateObject("roRegistrySection", "StreamVerse")
    if sec.Exists("auth_token")
        return sec.Read("auth_token")
    end if
    return invalid
end function

