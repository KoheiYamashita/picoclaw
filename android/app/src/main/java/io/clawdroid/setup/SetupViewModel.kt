package io.clawdroid.setup

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import io.clawdroid.backend.api.GatewaySettings
import io.clawdroid.backend.api.GatewaySettingsStore
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import kotlinx.serialization.json.JsonPrimitive
import kotlinx.serialization.json.buildJsonObject
import java.util.UUID

data class SetupUiState(
    val currentStep: Int = 0,
    val loading: Boolean = false,
    val error: String? = null,
    // Step 1: Gateway
    val gatewayPort: String = "18790",
    val gatewayApiKey: String = "",
    val step1Done: Boolean = false,
    // Step 2: LLM
    val llmModel: String = "",
    val llmApiKey: String = "",
    val llmBaseUrl: String = "",
    val step2Skipped: Boolean = false,
    // Step 3: Workspace
    val workspace: String = "",
    val dataDir: String = "",
    val step3Skipped: Boolean = false,
    // Step 4: WS + Agent
    val wsHost: String = "127.0.0.1",
    val wsPort: String = "18793",
    val wsPath: String = "/ws",
    val wsApiKey: String = "",
    val maxTokens: String = "8192",
    val contextWindow: String = "128000",
    val temperature: String = "0",
    val maxToolIterations: String = "10",
    val step4Skipped: Boolean = false,
) {
    val gatewayPortError: String?
        get() {
            if (gatewayPort.isEmpty()) return null
            val port = gatewayPort.toIntOrNull() ?: return "Invalid number"
            return if (port !in 1..65535) "1-65535" else null
        }

    val canProceedStep1: Boolean
        get() = gatewayPort.isNotEmpty() && gatewayPortError == null && gatewayApiKey.isNotEmpty()
}

class SetupViewModel(
    private val setupApiClient: SetupApiClient,
    private val settingsStore: GatewaySettingsStore,
) : ViewModel() {

    private val _uiState = MutableStateFlow(SetupUiState())
    val uiState: StateFlow<SetupUiState> = _uiState.asStateFlow()

    fun onGatewayPortChange(value: String) {
        if (value.isEmpty() || value.toIntOrNull() != null) {
            _uiState.update { it.copy(gatewayPort = value, error = null) }
        }
    }

    fun onGatewayApiKeyChange(value: String) {
        _uiState.update { it.copy(gatewayApiKey = value, error = null) }
    }

    fun generateApiKey() {
        _uiState.update { it.copy(gatewayApiKey = UUID.randomUUID().toString(), error = null) }
    }

    fun onLlmModelChange(value: String) = _uiState.update { it.copy(llmModel = value) }
    fun onLlmApiKeyChange(value: String) = _uiState.update { it.copy(llmApiKey = value) }
    fun onLlmBaseUrlChange(value: String) = _uiState.update { it.copy(llmBaseUrl = value) }

    fun onWorkspaceChange(value: String) = _uiState.update { it.copy(workspace = value) }
    fun onDataDirChange(value: String) = _uiState.update { it.copy(dataDir = value) }

    fun onWsHostChange(value: String) = _uiState.update { it.copy(wsHost = value) }
    fun onWsPortChange(value: String) {
        if (value.isEmpty() || value.toIntOrNull() != null) {
            _uiState.update { it.copy(wsPort = value) }
        }
    }
    fun onWsPathChange(value: String) = _uiState.update { it.copy(wsPath = value) }
    fun onWsApiKeyChange(value: String) = _uiState.update { it.copy(wsApiKey = value) }
    fun onMaxTokensChange(value: String) {
        if (value.isEmpty() || value.toIntOrNull() != null) {
            _uiState.update { it.copy(maxTokens = value) }
        }
    }
    fun onContextWindowChange(value: String) {
        if (value.isEmpty() || value.toIntOrNull() != null) {
            _uiState.update { it.copy(contextWindow = value) }
        }
    }
    fun onTemperatureChange(value: String) {
        if (value.isEmpty() || value.toDoubleOrNull() != null) {
            _uiState.update { it.copy(temperature = value) }
        }
    }
    fun onMaxToolIterationsChange(value: String) {
        if (value.isEmpty() || value.toIntOrNull() != null) {
            _uiState.update { it.copy(maxToolIterations = value) }
        }
    }

    fun submitInit() {
        viewModelScope.launch {
            val state = _uiState.value
            if (!state.canProceedStep1 || state.loading) return@launch

            _uiState.update { it.copy(loading = true, error = null) }

            val port = state.gatewayPort.toIntOrNull() ?: 18790
            val body = buildJsonObject {
                put("gateway", buildJsonObject {
                    put("port", JsonPrimitive(port))
                    put("api_key", JsonPrimitive(state.gatewayApiKey))
                })
            }

            try {
                setupApiClient.init(body)
                // Persist gateway settings locally
                settingsStore.update(GatewaySettings(httpPort = port, apiKey = state.gatewayApiKey))
                _uiState.update { it.copy(loading = false, step1Done = true, currentStep = 1) }
            } catch (e: Exception) {
                _uiState.update { it.copy(loading = false, error = e.message ?: "Init failed") }
            }
        }
    }

    fun skipStep(step: Int) {
        _uiState.update {
            when (step) {
                2 -> it.copy(step2Skipped = true, currentStep = 2)
                3 -> it.copy(step3Skipped = true, currentStep = 3)
                4 -> it.copy(step4Skipped = true, currentStep = 4)
                else -> it
            }
        }
    }

    fun nextStep(step: Int) {
        _uiState.update { it.copy(currentStep = step) }
    }

    fun submitComplete(onComplete: () -> Unit) {
        viewModelScope.launch {
            val state = _uiState.value
            if (state.loading) return@launch

            _uiState.update { it.copy(loading = true, error = null) }

            val body = buildJsonObject {
                if (!state.step2Skipped) {
                    put("llm", buildJsonObject {
                        if (state.llmModel.isNotBlank()) put("model", JsonPrimitive(state.llmModel))
                        if (state.llmApiKey.isNotBlank()) put("api_key", JsonPrimitive(state.llmApiKey))
                        if (state.llmBaseUrl.isNotBlank()) put("base_url", JsonPrimitive(state.llmBaseUrl))
                    })
                }
                if (!state.step3Skipped) {
                    put("agents", buildJsonObject {
                        put("defaults", buildJsonObject {
                            if (state.workspace.isNotBlank()) put("workspace", JsonPrimitive(state.workspace))
                            if (state.dataDir.isNotBlank()) put("data_dir", JsonPrimitive(state.dataDir))
                        })
                    })
                }
                if (!state.step4Skipped) {
                    put("channels", buildJsonObject {
                        put("websocket", buildJsonObject {
                            if (state.wsHost.isNotBlank()) put("host", JsonPrimitive(state.wsHost))
                            state.wsPort.toIntOrNull()?.let { put("port", JsonPrimitive(it)) }
                            if (state.wsPath.isNotBlank()) put("path", JsonPrimitive(state.wsPath))
                            if (state.wsApiKey.isNotBlank()) put("api_key", JsonPrimitive(state.wsApiKey))
                        })
                    })
                    put("agents_extra", buildJsonObject {
                        put("defaults", buildJsonObject {
                            state.maxTokens.toIntOrNull()?.let { put("max_tokens", JsonPrimitive(it)) }
                            state.contextWindow.toIntOrNull()?.let { put("context_window", JsonPrimitive(it)) }
                            state.temperature.toDoubleOrNull()?.let { put("temperature", JsonPrimitive(it)) }
                            state.maxToolIterations.toIntOrNull()?.let { put("max_tool_iterations", JsonPrimitive(it)) }
                        })
                    })
                }
            }

            try {
                setupApiClient.complete(body)
                _uiState.update { it.copy(loading = false) }
                onComplete()
            } catch (e: Exception) {
                _uiState.update { it.copy(loading = false, error = e.message ?: "Complete failed") }
            }
        }
    }
}
