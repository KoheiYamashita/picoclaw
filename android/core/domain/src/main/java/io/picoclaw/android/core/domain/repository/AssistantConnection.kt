package io.picoclaw.android.core.domain.repository

import io.picoclaw.android.core.domain.model.AssistantMessage
import io.picoclaw.android.core.domain.model.ConnectionState
import kotlinx.coroutines.flow.SharedFlow
import kotlinx.coroutines.flow.StateFlow

interface AssistantConnection {
    val messages: SharedFlow<AssistantMessage>
    val statusText: StateFlow<String?>
    val connectionState: StateFlow<ConnectionState>

    fun connect(wsUrl: String)
    fun disconnect()
    suspend fun send(text: String, images: List<String> = emptyList(), inputMode: String = "assistant")
}
