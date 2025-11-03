package com.streamverse.mobile.ui.content.detail

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import coil.compose.AsyncImage
import com.streamverse.mobile.models.Content

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ContentDetailScreen(
    content: Content,
    onDismiss: () -> Unit,
    onPlay: () -> Unit
) {
    AlertDialog(
        onDismissRequest = onDismiss,
        modifier = Modifier.fillMaxWidth()
    ) {
        Column(
            modifier = Modifier
                .verticalScroll(rememberScrollState())
                .padding(16.dp)
        ) {
            AsyncImage(
                model = content.backdropUrl,
                contentDescription = content.title,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(200.dp)
            )
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Text(
                text = content.title,
                style = MaterialTheme.typography.headlineMedium
            )
            
            Spacer(modifier = Modifier.height(8.dp))
            
            Row {
                Text("${content.releaseYear} • ")
                Text("${String.format("%.1f", content.rating)} ⭐ • ")
                Text(content.genre)
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Text(
                text = content.description,
                style = MaterialTheme.typography.bodyMedium
            )
            
            Spacer(modifier = Modifier.height(24.dp))
            
            Button(
                onClick = onPlay,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text("Play")
            }
        }
    }
}

