# Testing Guide - Android TV App

## Running Tests

### Unit Tests

```bash
# Run all unit tests
./gradlew test

# Run specific test class
./gradlew test --tests "com.streamverse.tv.viewmodel.MainViewModelTest"

# Generate test coverage report
./gradlew testDebugUnitTest jacocoTestReport
```

### Integration Tests

```bash
# Run on connected Android TV device/emulator
./gradlew connectedAndroidTest

# Run specific test
./gradlew connectedAndroidTest --tests "com.streamverse.tv.MainActivityTest"
```

## Test Structure

```
app/src/
├── test/              # Unit tests (run on JVM)
│   └── java/com/streamverse/tv/
│       ├── viewmodel/
│       │   ├── MainViewModelTest.kt
│       │   ├── AuthViewModelTest.kt
│       │   └── SearchViewModelTest.kt
│       └── data/repository/
│           └── ContentRepositoryTest.kt
│
└── androidTest/       # Instrumented tests (run on device)
    └── java/com/streamverse/tv/
        ├── MainActivityTest.kt
        └── PlaybackVideoActivityTest.kt
```

## Test Coverage

Current coverage:
- ✅ ViewModel unit tests
- ✅ Repository unit tests
- ⏳ Fragment/Activity integration tests (pending)

## Writing Tests

### Example: ViewModel Test

```kotlin
@Test
fun `test viewModel behavior`() = runTest {
    // Given
    val mockRepository = mock<ContentRepository>()
    whenever(mockRepository.getHomeContentRows()).thenReturn(emptyList())
    
    // When
    val viewModel = MainViewModel(mockRepository)
    viewModel.loadContent()
    
    // Then
    assertEquals(emptyList(), viewModel.contentRows.value)
}
```

## Mocking

We use Mockito-Kotlin for mocking:
- Mock repositories
- Mock API services
- Verify method calls

## Test Data

Create test data factories for consistent test data:
- `ContentFactory` - Create test Content objects
- `AuthResponseFactory` - Create test auth responses

