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
import androidx.compose.foundation.text.KeyboardOptions
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
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import io.clawdroid.core.ui.theme.DeepBlack
import io.clawdroid.core.ui.theme.NeonCyan
import io.clawdroid.core.ui.theme.TextPrimary
import io.clawdroid.core.ui.theme.TextSecondary

@Composable
fun SetupStep4ChatScreen(viewModel: SetupViewModel) {
    val uiState by viewModel.uiState.collectAsState()
    var wsApiKeyHidden by remember { mutableStateOf(true) }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp)
            .verticalScroll(rememberScrollState()),
        verticalArrangement = Arrangement.spacedBy(12.dp),
    ) {
        Spacer(Modifier.height(32.dp))

        Text("Step 4 of 5", style = MaterialTheme.typography.labelMedium, color = TextSecondary)
        Text("WebSocket & Agent", style = MaterialTheme.typography.headlineMedium, color = TextPrimary)
        Text(
            "Configure the WebSocket channel and agent parameters.",
            style = MaterialTheme.typography.bodyMedium,
            color = TextSecondary,
        )

        Spacer(Modifier.height(4.dp))

        // WebSocket section
        Text("WebSocket", style = MaterialTheme.typography.titleSmall, color = NeonCyan)

        OutlinedTextField(
            value = uiState.wsHost,
            onValueChange = viewModel::onWsHostChange,
            label = { Text("Host", color = TextSecondary) },
            singleLine = true,
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.wsPort,
            onValueChange = viewModel::onWsPortChange,
            label = { Text("Port", color = TextSecondary) },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.wsPath,
            onValueChange = viewModel::onWsPathChange,
            label = { Text("Path", color = TextSecondary) },
            singleLine = true,
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.wsApiKey,
            onValueChange = viewModel::onWsApiKeyChange,
            label = { Text("WS API Key", color = TextSecondary) },
            singleLine = true,
            visualTransformation = if (wsApiKeyHidden) PasswordVisualTransformation() else VisualTransformation.None,
            trailingIcon = {
                TextButton(onClick = { wsApiKeyHidden = !wsApiKeyHidden }) {
                    Text(
                        if (wsApiKeyHidden) "Show" else "Hide",
                        color = NeonCyan,
                        style = MaterialTheme.typography.labelSmall,
                    )
                }
            },
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        Spacer(Modifier.height(8.dp))

        // Agent section
        Text("Agent Defaults", style = MaterialTheme.typography.titleSmall, color = NeonCyan)

        OutlinedTextField(
            value = uiState.maxTokens,
            onValueChange = viewModel::onMaxTokensChange,
            label = { Text("Max Tokens", color = TextSecondary) },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.contextWindow,
            onValueChange = viewModel::onContextWindowChange,
            label = { Text("Context Window", color = TextSecondary) },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.temperature,
            onValueChange = viewModel::onTemperatureChange,
            label = { Text("Temperature", color = TextSecondary) },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        OutlinedTextField(
            value = uiState.maxToolIterations,
            onValueChange = viewModel::onMaxToolIterationsChange,
            label = { Text("Max Tool Iterations", color = TextSecondary) },
            singleLine = true,
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
            colors = setupFieldColors(),
            modifier = Modifier.fillMaxWidth(),
        )

        Spacer(Modifier.weight(1f))

        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
        ) {
            TextButton(onClick = { viewModel.skipStep(4) }) {
                Text("Set up later", color = TextSecondary)
            }
            Button(
                onClick = { viewModel.nextStep(4) },
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
