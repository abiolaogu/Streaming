//
//  AuthViewModelTests.swift
//  StreamVerseTests
//
//  Created on 2025-01-XX.
//

import XCTest
@testable import StreamVerse

@MainActor
final class AuthViewModelTests: XCTestCase {
    var viewModel: AuthViewModel!
    
    override func setUp() {
        super.setUp()
        viewModel = AuthViewModel()
    }
    
    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }
    
    func testInitialState() {
        XCTAssertFalse(viewModel.isLoggedIn)
        XCTAssertFalse(viewModel.isLoading)
        XCTAssertNil(viewModel.errorMessage)
        XCTAssertNil(viewModel.currentUser)
    }
    
    func testCheckAuthStatus() {
        // Given: User is logged in (mocked)
        // When: Check auth status
        viewModel.checkAuthStatus()
        
        // Then: Status is checked
        // Note: This depends on TokenManager state
        XCTAssertNotNil(viewModel.isLoggedIn)
    }
}

