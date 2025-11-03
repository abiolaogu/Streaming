package com.streamverse.mobile.ui.content

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.ui.Alignment
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.compose.ui.platform.LocalContext
import coil.compose.AsyncImage
import com.streamverse.mobile.data.repository.ContentRepository
import com.streamverse.mobile.models.Content
import com.streamverse.mobile.models.ContentRow
import com.streamverse.mobile.ui.content.detail.ContentDetailScreen
import com.streamverse.mobile.ui.player.VideoPlayerScreen
import com.streamverse.mobile.viewmodel.AuthViewModel
import com.streamverse.mobile.viewmodel.ContentViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ContentListScreen(
    authViewModel: AuthViewModel,
    onLogout: () -> Unit
) {
    val context = LocalContext.current
    val contentViewModel: ContentViewModel = viewModel(
        factory = ContentViewModel.Factory(
            ContentRepository.create(context)
        )
    )

    var selectedContent by remember { mutableStateOf<Content?>(null) }
    var showPlayer by remember { mutableStateOf(false) }
    var searchQuery by remember { mutableStateOf("") }
    var searchResults by remember { mutableStateOf<List<Content>>(emptyList()) }

    val contentRows by contentViewModel.contentRows.collectAsState()
    val isLoading by contentViewModel.isLoading.collectAsState()
    val errorMessage by contentViewModel.errorMessage.collectAsState()

    LaunchedEffect(Unit) {
        contentViewModel.loadHomeContent()
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("StreamVerse") },
                actions = {
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Default.AccountCircle, contentDescription = "Logout")
                    }
                }
            )
        }
    ) { padding ->
        when {
            isLoading && contentRows.isEmpty() -> {
                Box(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding),
                    contentAlignment = androidx.compose.ui.Alignment.Center
                ) {
                    CircularProgressIndicator()
                }
            }
            errorMessage != null -> {
                Column(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding),
                    horizontalAlignment = androidx.compose.ui.Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    Text("Error: $errorMessage")
                    Button(onClick = { contentViewModel.loadHomeContent() }) {
                        Text("Retry")
                    }
                }
            }
            else -> {
                LazyColumn(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(padding),
                    verticalArrangement = Arrangement.spacedBy(24.dp)
                ) {
                    items(contentRows) { row ->
                        ContentRowSection(
                            row = row,
                            onContentClick = { content ->
                                selectedContent = content
                            }
                        )
                    }
                }
            }
        }
    }

    selectedContent?.let { content ->
        ContentDetailScreen(
            content = content,
            onDismiss = { selectedContent = null },
            onPlay = { showPlayer = true }
        )
    }

    if (showPlayer && selectedContent != null) {
        VideoPlayerScreen(
            content = selectedContent!!,
            onDismiss = { showPlayer = false }
        )
    }
}

@Composable
fun ContentRowSection(
    row: ContentRow,
    onContentClick: (Content) -> Unit
) {
    Column(
        modifier = Modifier.padding(horizontal = 16.dp)
    ) {
        Text(
            text = row.title,
            style = MaterialTheme.typography.titleLarge,
            modifier = Modifier.padding(bottom = 12.dp)
        )
        LazyRow(
            horizontalArrangement = Arrangement.spacedBy(12.dp)
        ) {
            items(row.items) { content ->
                ContentCard(
                    content = content,
                    onClick = { onContentClick(content) }
                )
            }
        }
    }
}

@Composable
fun ContentCard(
    content: Content,
    onClick: () -> Unit
) {
    Card(
        onClick = onClick,
        modifier = Modifier
            .width(150.dp)
            .height(225.dp)
    ) {
        Column {
            AsyncImage(
                model = content.posterUrl,
                contentDescription = content.title,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(200.dp)
            )
            Text(
                text = content.title,
                style = MaterialTheme.typography.bodySmall,
                maxLines = 2,
                modifier = Modifier.padding(8.dp)
            )
        }
    }
}

