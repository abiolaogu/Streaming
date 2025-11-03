package com.streamverse.tv

import android.app.Application
import org.koin.android.ext.koin.androidContext
import org.koin.android.ext.koin.androidLogger
import org.koin.core.context.startKoin
import org.koin.core.logger.Level
import org.koin.dsl.module

/**
 * Application class for StreamVerse Android TV app.
 * Initializes dependency injection with Koin.
 */
class StreamVerseApplication : Application() {

    override fun onCreate() {
        super.onCreate()
        
        startKoin {
            androidLogger(Level.INFO)
            androidContext(this@StreamVerseApplication)
            modules(appModule)
        }
    }
}

// Koin dependency injection module
val appModule = module {
    // Network components
    single { createOkHttpClient() }
    single { createRetrofit(get()) }
    single<com.streamverse.tv.data.api.ContentApiService> { get<retrofit2.Retrofit>().create() }
    
    // Repository
    single { com.streamverse.tv.data.repository.ContentRepository(get()) }
    
    // ViewModels
    viewModel { com.streamverse.tv.viewmodel.MainViewModel(get()) }
    viewModel { com.streamverse.tv.viewmodel.BrowseViewModel(get()) }
}

fun createOkHttpClient(): okhttp3.OkHttpClient {
    return okhttp3.OkHttpClient.Builder()
        .connectTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
        .readTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
        .writeTimeout(30, java.util.concurrent.TimeUnit.SECONDS)
        .addInterceptor { chain ->
            val request = chain.request().newBuilder()
                .addHeader("Accept", "application/json")
                .addHeader("Content-Type", "application/json")
                .build()
            chain.proceed(request)
        }
        .build()
}

fun createRetrofit(okHttpClient: okhttp3.OkHttpClient): retrofit2.Retrofit {
    return retrofit2.Retrofit.Builder()
        .baseUrl("https://api.streamverse.com/") // TODO: Move to BuildConfig or config file
        .client(okHttpClient)
        .addConverterFactory(retrofit2.converter.gson.GsonConverterFactory.create())
        .build()
}

