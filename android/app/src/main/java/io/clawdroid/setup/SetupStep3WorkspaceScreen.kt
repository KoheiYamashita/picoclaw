package io.clawdroid.setup

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.material3.TextButton
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import io.clawdroid.core.ui.theme.DeepBlack
import io.clawdroid.core.ui.theme.NeonCyan
import io.clawdroid.core.ui.theme.TextPrimary
import io.clawdroid.core.ui.theme.TextSecondary

@Composable
fun SetupStep3WorkspaceScreen(viewModel: SetupViewModel) {
    val uiState by viewModel.uiState.collectAsState()

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp)
            .verticalScroll(rememberScrollState()),
        verticalArrangement = Arrangement.spacedBy(16.dp),
    ) {
        Spacer(Modifier.height(32.dp))

        Text("Step 3 of 5", style = MaterialTheme.typography.labelMedium, color = TextSecondary)
        Text("Workspace & Data", style = MaterialTheme.typography.headlineMedium, color = TextPrimary)
        Text(
            "Set the workspace and data directories used by the agent.",
            style = MaterialTheme.typography.bodyMedium,
            color = TextSecondary,
        )

        Spacer(Modifier.height(8.dp))

        OutlinedTextField(
            value = uiState.workspace,
            onValueChange = viewModel::onWorkspaceChange,
            label = { Text("Workspace", color = TextSecondary) },
            placeholder = { Text("~/.clawdroid/workspace", color = TextSecondary.copy(alpha = 0.5f)) },
            singleLine = true,
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.dataDir,
            onValueChange = viewModel::onDataDirChange,
            label = { Text("Data Directory", color = TextSecondary) },
            placeholder = { Text("~/.clawdroid/data", color = TextSecondary.copy(alpha = 0.5f)) },
            singleLine = true,
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        Spacer(Modifier.weight(1f))

        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
        ) {
            TextButton(onClick = { viewModel.skipStep(3) }) {
                Text("Set up later", color = TextSecondary)
            }
            Button(
                onClick = { viewModel.nextStep(3) },
                colors = ButtonDefaults.buttonColors(
                    containerColor = NeonCyan,
                    contentColor = DeepBlack,
                ),
            ) {
                Text("Next")
            }
        }
    }
}
