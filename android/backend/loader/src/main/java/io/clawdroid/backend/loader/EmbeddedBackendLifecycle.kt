package io.clawdroid.backend.loader

import android.content.Context
import android.content.Intent
import io.clawdroid.backend.api.BackendLifecycle
import io.clawdroid.backend.api.BackendState
import kotlinx.coroutines.flow.StateFlow

class EmbeddedBackendLifecycle(
    private val context: Context,
    private val processManager: GatewayProcessManager,
) : BackendLifecycle {

    override val isManaged: Boolean = true

    override val state: StateFlow<BackendState> = processManager.state

    override suspend fun start() {
        val intent = Intent(context, GatewayService::class.java)
        context.startForegroundService(intent)
    }

    override suspend fun stop() {
        context.stopService(Intent(context, GatewayService::class.java))
    }
}
