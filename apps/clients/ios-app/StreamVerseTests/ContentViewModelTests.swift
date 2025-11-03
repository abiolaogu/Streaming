//
//  ContentViewModelTests.swift
//  StreamVerseTests
//
//  Created on 2025-01-XX.
//

import XCTest
@testable import StreamVerse

@MainActor
final class ContentViewModelTests: XCTestCase {
    var viewModel: ContentViewModel!
    
    override func setUp() {
        super.setUp()
        viewModel = ContentViewModel()
    }
    
    override func tearDown() {
        viewModel = nil
        super.tearDown()
    }
    
    func testInitialState() {
        XCTAssertTrue(viewModel.contentRows.isEmpty)
        XCTAssertFalse(viewModel.isLoading)
        XCTAssertNil(viewModel.errorMessage)
    }
    
    func testSearchWithShortQuery() async {
        // When: Search with query less than 2 characters
        let results = await viewModel.search(query: "a")
        
        // Then: Returns empty array
        XCTAssertTrue(results.isEmpty)
    }
}

