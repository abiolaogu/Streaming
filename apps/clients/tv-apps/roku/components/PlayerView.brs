'
' Video Player View Controller
'

function Init()
    m.videoPlayer = m.top.findNode("videoPlayer")
    m.top.ObserveField("content", "OnContentSet")
    
    ' Set up video event listener
    m.videoPlayer.ObserveField("state", "OnVideoStateChanged")
end function

function OnContentSet()
    content = m.top.content
    
    ' Create video content node
    videoContent = CreateObject("roSGNode", "ContentNode")
    videoContent.url = content.streamUrl
    videoContent.title = content.title
    
    ' Configure DRM if needed
    if content.isDrmProtected
        drmConfig = CreateObject("roAssociativeArray")
        
        if content.drmType = "widevine"
            drmConfig.type = "Widevine"
        else if content.drmType = "playready"
            drmConfig.type = "PlayReady"
        end if
        
        drmConfig.licenseServerURL = "https://drm.streamverse.com/v1/" + LCase(content.drmType) + "/license"
        
        ' Add auth token to license request
        token = GetAuthToken()
        if token <> invalid
            drmConfig.customData = FormatJson({
                Authorization: "Bearer " + token
                "X-Content-ID": content.id
            })
        end if
        
        videoContent.drmParams = drmConfig
    end if
    
    ' Set content and play
    m.videoPlayer.content = videoContent
    m.videoPlayer.control = "play"
end function

function OnVideoStateChanged()
    state = m.videoPlayer.state
    
    if state = "finished"
        ' Return to detail view
        m.top.getParent().CallFunc("ShowDetail")
    else if state = "error"
        ' Show error message
        ShowErrorDialog("Playback error occurred")
    end if
end function

function GetAuthToken() as String
    sec = CreateObject("roRegistrySection", "StreamVerse")
    if sec.Exists("auth_token")
        return sec.Read("auth_token")
    end if
    return invalid
end function

function ShowErrorDialog(message as String)
    dialog = CreateObject("roSGNode", "StandardMessageDialog")
    dialog.title = "Error"
    dialog.message = message
    dialog.buttons = ["OK"]
    m.top.getParent().dialog = dialog
end function

