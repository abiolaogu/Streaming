//
//  DRMHelper.swift
//  StreamVerse
//
//  Created on 2025-01-XX.
//

import AVFoundation

/// Helper for DRM configuration (FairPlay).
class DRMHelper {
    static let shared = DRMHelper()
    
    private let licenseServerURL: String
    
    private init() {
        // Load from Info.plist or Build Settings
        self.licenseServerURL = Bundle.main.object(forInfoDictionaryKey: "DRM_LICENSE_SERVER") as? String
            ?? "https://drm.streamverse.com/v1/fairplay/license"
    }
    
    /// Configure DRM for AVAsset.
    func configureDRM(for asset: AVURLAsset, content: Content, drmType: String) {
        guard drmType.lowercased() == "fairplay" else {
            print("Unsupported DRM type: \(drmType)")
            return
        }
        
        // FairPlay configuration
        asset.resourceLoader.setDelegate(self, queue: DispatchQueue.main)
    }
}

extension DRMHelper: AVAssetResourceLoaderDelegate {
    func resourceLoader(
        _ resourceLoader: AVAssetResourceLoader,
        shouldWaitForLoadingOfRequestedResource loadingRequest: AVAssetResourceLoadingRequest
    ) -> Bool {
        guard let url = loadingRequest.request.url,
              url.scheme == "skd" else {
            return false
        }
        
        // Handle FairPlay license request
        handleFairPlayLicenseRequest(loadingRequest)
        
        return true
    }
    
    private func handleFairPlayLicenseRequest(_ request: AVAssetResourceLoadingRequest) {
        guard let contentId = request.request.url?.host,
              let dataRequest = request.dataRequest else {
            request.finishLoading(with: NSError(domain: "DRM", code: -1))
            return
        }
        
        // Get auth token
        let token = TokenManager.shared.getAccessToken() ?? ""
        
        // Create license request
        var licenseRequest = URLRequest(url: URL(string: licenseServerURL)!)
        licenseRequest.httpMethod = "POST"
        licenseRequest.setValue("Bearer \(token)", forHTTPHeaderField: "Authorization")
        licenseRequest.setValue(contentId, forHTTPHeaderField: "X-Content-ID")
        licenseRequest.setValue("application/octet-stream", forHTTPHeaderField: "Content-Type")
        
        // Use the request's data as the body
        if let requestData = dataRequest.requestedOffset == 0 ? dataRequest.requestedData : nil {
            licenseRequest.httpBody = requestData
        }
        
        // Perform license request
        URLSession.shared.dataTask(with: licenseRequest) { data, response, error in
            if let error = error {
                request.finishLoading(with: error)
                return
            }
            
            guard let data = data,
                  let response = response as? HTTPURLResponse,
                  (200...299).contains(response.statusCode) else {
                request.finishLoading(with: NSError(domain: "DRM", code: -1))
                return
            }
            
            // Provide license data
            dataRequest.respond(with: data)
            request.finishLoading()
        }.resume()
    }
}

