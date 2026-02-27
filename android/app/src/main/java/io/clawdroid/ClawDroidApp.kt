package io.clawdroid

import android.app.Application
import io.clawdroid.backend.api.GatewaySettingsStore
import io.clawdroid.backend.config.configModule
import io.clawdroid.core.data.remote.WebSocketClient
import io.clawdroid.di.appModule
import io.clawdroid.receiver.NotificationHelper
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.launch
import org.koin.android.ext.koin.androidContext
import org.koin.core.context.startKoin

class ClawDroidApp : Application() {
    override fun onCreate() {
        super.onCreate()
        val koinApp = startKoin {
            androidContext(this@ClawDroidApp)
            modules(appModule, configModule)
        }
        NotificationHelper.createNotificationChannel(this)

        val koin = koinApp.koin
        val settingsStore: GatewaySettingsStore = koin.get()
        val wsClient: WebSocketClient = koin.get()
        val scope: CoroutineScope = koin.get()
        scope.launch {
            settingsStore.settings.collect { s ->
                if (wsClient.wsUrl != s.wsUrl || wsClient.apiKey != s.apiKey) {
                    wsClient.disconnect()
                    wsClient.wsUrl = s.wsUrl
                    wsClient.apiKey = s.apiKey
                    wsClient.connect()
                }
            }
        }
    }
}
