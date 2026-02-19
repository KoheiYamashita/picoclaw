// TTS settings are currently used only within the chat feature, so this screen
// is placed under feature/chat. If TTS settings become shared across multiple
// features in the future, consider extracting them into a dedicated feature/settings module.
package io.picoclaw.android.feature.chat.screen

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material3.Button
import androidx.compose.material3.DropdownMenuItem
import androidx.compose.material3.ExperimentalMaterial3Api
import androidx.compose.material3.ExposedDropdownMenuAnchorType
import androidx.compose.material3.ExposedDropdownMenuBox
import androidx.compose.material3.ExposedDropdownMenuDefaults
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Slider
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableFloatStateOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import io.picoclaw.android.core.domain.model.TtsEngineInfo
import io.picoclaw.android.core.domain.model.TtsVoiceInfo
import io.picoclaw.android.feature.chat.SettingsViewModel
import org.koin.androidx.compose.koinViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SettingsScreen(
    onNavigateBack: () -> Unit,
    viewModel: SettingsViewModel = koinViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Settings") },
                navigationIcon = {
                    IconButton(onClick = onNavigateBack) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(24.dp)
        ) {
            Text("Text-to-Speech", style = MaterialTheme.typography.titleMedium)

            EngineSelector(
                selectedEngine = uiState.ttsConfig.enginePackageName,
                engines = uiState.availableEngines,
                onEngineSelected = viewModel::onEngineSelected
            )

            VoiceSelector(
                selectedVoiceName = uiState.ttsConfig.voiceName,
                voices = uiState.availableVoices,
                onVoiceSelected = viewModel::onVoiceSelected
            )

            SliderSetting(
                label = "Speed",
                value = uiState.ttsConfig.speechRate,
                valueRange = 0.5f..2.0f,
                onValueChangeFinished = viewModel::onSpeechRateChanged
            )

            SliderSetting(
                label = "Pitch",
                value = uiState.ttsConfig.pitch,
                valueRange = 0.5f..2.0f,
                onValueChangeFinished = viewModel::onPitchChanged
            )

            Button(
                onClick = viewModel::onTestSpeak,
                enabled = !uiState.isTesting
            ) {
                Text(if (uiState.isTesting) "Speaking..." else "Test Voice")
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun EngineSelector(
    selectedEngine: String?,
    engines: List<TtsEngineInfo>,
    onEngineSelected: (String?) -> Unit
) {
    var expanded by remember { mutableStateOf(false) }
    val displayText = if (selectedEngine == null) {
        "System Default"
    } else {
        engines.find { it.packageName == selectedEngine }?.label ?: selectedEngine
    }

    ExposedDropdownMenuBox(
        expanded = expanded,
        onExpandedChange = { expanded = it }
    ) {
        OutlinedTextField(
            value = displayText,
            onValueChange = {},
            readOnly = true,
            label = { Text("Engine") },
            trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded) },
            modifier = Modifier
                .menuAnchor(ExposedDropdownMenuAnchorType.PrimaryNotEditable)
                .fillMaxWidth()
        )
        ExposedDropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false }
        ) {
            DropdownMenuItem(
                text = { Text("System Default") },
                onClick = {
                    onEngineSelected(null)
                    expanded = false
                }
            )
            engines.forEach { engine ->
                DropdownMenuItem(
                    text = { Text(engine.label) },
                    onClick = {
                        onEngineSelected(engine.packageName)
                        expanded = false
                    }
                )
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
private fun VoiceSelector(
    selectedVoiceName: String?,
    voices: List<TtsVoiceInfo>,
    onVoiceSelected: (String?) -> Unit
) {
    var expanded by remember { mutableStateOf(false) }
    val displayText = if (selectedVoiceName == null) {
        "System Default"
    } else {
        voices.find { it.name == selectedVoiceName }?.displayLabel ?: selectedVoiceName
    }

    ExposedDropdownMenuBox(
        expanded = expanded,
        onExpandedChange = { expanded = it }
    ) {
        OutlinedTextField(
            value = displayText,
            onValueChange = {},
            readOnly = true,
            label = { Text("Voice") },
            trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded) },
            modifier = Modifier
                .menuAnchor(ExposedDropdownMenuAnchorType.PrimaryNotEditable)
                .fillMaxWidth()
        )
        ExposedDropdownMenu(
            expanded = expanded,
            onDismissRequest = { expanded = false }
        ) {
            DropdownMenuItem(
                text = { Text("System Default") },
                onClick = {
                    onVoiceSelected(null)
                    expanded = false
                }
            )
            voices.forEach { voice ->
                DropdownMenuItem(
                    text = { Text(voice.displayLabel) },
                    onClick = {
                        onVoiceSelected(voice.name)
                        expanded = false
                    }
                )
            }
        }
    }
}

@Composable
private fun SliderSetting(
    label: String,
    value: Float,
    valueRange: ClosedFloatingPointRange<Float>,
    onValueChangeFinished: (Float) -> Unit
) {
    var localValue by remember(value) { mutableFloatStateOf(value) }

    Column {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(label, style = MaterialTheme.typography.bodyLarge)
            Text(
                "%.1f".format(localValue),
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
            )
        }
        Slider(
            value = localValue,
            onValueChange = { localValue = it },
            onValueChangeFinished = { onValueChangeFinished(localValue) },
            valueRange = valueRange,
            steps = 14
        )
    }
}
