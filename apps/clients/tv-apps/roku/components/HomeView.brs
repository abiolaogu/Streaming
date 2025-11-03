'
' Home View Controller
'

function Init()
    m.contentRows = m.top.findNode("contentRows")
    
    ' Configure RowList
    m.contentRows.content = CreateObject("roSGNode", "ContentNode")
    
    ' Load home content
    LoadContent()
end function

function LoadContent()
    ' Make API request
    request = CreateObject("roUrlTransfer")
    request.SetUrl("https://api.streamverse.com/api/v1/content/home")
    
    ' Add auth token
    token = GetAuthToken()
    if token <> invalid
        request.AddHeader("Authorization", "Bearer " + token)
    end if
    
    request.SetCertificatesFile("common:/certs/ca-bundle.crt")
    request.InitClientCertificates()
    
    response = request.GetToString()
    
    if response <> ""
        json = ParseJson(response)
        if json <> invalid
            PopulateContentRows(json)
        end if
    end if
end function

function PopulateContentRows(rows)
    for each row in rows
        rowNode = CreateObject("roSGNode", "ContentNode")
        rowNode.title = row.title
        
        for each item in row.items
            itemNode = rowNode.CreateChild("ContentNode")
            itemNode.title = item.title
            itemNode.description = item.description
            itemNode.HDPosterUrl = item.posterUrl
            itemNode.SDPosterUrl = item.posterUrl
            itemNode.contentId = item.id
            itemNode.streamUrl = item.streamUrl
            itemNode.isDrmProtected = item.isDrmProtected
            itemNode.drmType = item.drmType
        end for
        
        m.contentRows.appendChild(rowNode)
    end for
    
    m.contentRows.SetFocus(true)
end function

function OnItemSelected()
    selectedItem = m.contentRows.content.GetChild(m.contentRows.rowItemSelected[0]).GetChild(m.contentRows.rowItemSelected[1])
    m.top.contentSelected = {
        id: selectedItem.contentId
        title: selectedItem.title
        description: selectedItem.description
        streamUrl: selectedItem.streamUrl
        isDrmProtected: selectedItem.isDrmProtected
        drmType: selectedItem.drmType
    }
end function

function GetAuthToken() as String
    sec = CreateObject("roRegistrySection", "StreamVerse")
    if sec.Exists("auth_token")
        return sec.Read("auth_token")
    end if
    return invalid
end function

