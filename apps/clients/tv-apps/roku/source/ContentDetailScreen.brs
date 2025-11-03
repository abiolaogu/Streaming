sub init()
    m.backdrop = m.top.findNode("backdrop")
    m.title = m.top.findNode("title")
    m.meta = m.top.findNode("meta")
    m.description = m.top.findNode("description")
    m.playButton = m.top.findNode("playButton")
    
    m.playButton.setFocus(true)
    m.playButton.observeField("buttonSelected", "onPlayButtonSelected")
    
    if m.top.content <> invalid
        populateContent()
    end if
end sub

sub populateContent()
    content = m.top.content
    
    m.backdrop.uri = content.backdropUrl
    m.title.text = content.title
    m.meta.text = content.releaseYear.toStr() + " • " + content.rating.toStr() + " ⭐ • " + content.genre
    m.description.text = content.description
end sub

sub onPlayButtonSelected()
    content = m.top.content
    
    ' Create video player
    video = CreateObject("roSGNode", "Video")
    video.content = content
    video.control = "play"
    
    m.top.appendChild(video)
    video.setFocus(true)
end sub

