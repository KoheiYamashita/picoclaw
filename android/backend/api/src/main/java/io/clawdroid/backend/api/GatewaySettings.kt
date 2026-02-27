package io.clawdroid.backend.api

data class GatewaySettings(
    val wsPort: Int = 18793,
    val httpPort: Int = 18790,
    val apiKey: String = "",
) {
    val wsUrl: String get() = "ws://127.0.0.1:$wsPort/ws"
    val httpBaseUrl: String get() = "http://127.0.0.1:$httpPort"
}
