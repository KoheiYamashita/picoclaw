package io.clawdroid.settings

import androidx.lifecycle.SavedStateHandle
import io.clawdroid.backend.api.GatewaySettings
import io.clawdroid.backend.api.GatewaySettingsStore
import io.clawdroid.backend.config.ConfigApiClient
import io.clawdroid.backend.config.SaveConfigResult
import io.mockk.coEvery
import io.mockk.coVerify
import io.mockk.every
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.test.UnconfinedTestDispatcher
import kotlinx.coroutines.test.advanceUntilIdle
import kotlinx.coroutines.test.resetMain
import kotlinx.coroutines.test.runTest
import kotlinx.coroutines.test.setMain
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Assertions.assertEquals
import org.junit.jupiter.api.Assertions.assertFalse
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

@OptIn(ExperimentalCoroutinesApi::class)
class AppSettingsViewModelTest {

    private val testDispatcher = UnconfinedTestDispatcher()
    private val settingsFlow = MutableStateFlow(GatewaySettings(httpPort = 18790, apiKey = "test-key"))
    private val settingsStore = mockk<GatewaySettingsStore>(relaxed = true) {
        every { settings } returns settingsFlow
    }
    private val configApiClient = mockk<ConfigApiClient>(relaxed = true)
    private val savedStateHandle = SavedStateHandle()

    private lateinit var viewModel: AppSettingsViewModel

    @BeforeEach
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        viewModel = AppSettingsViewModel(savedStateHandle, settingsStore, configApiClient)
    }

    @AfterEach
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `init loads current settings`() {
        assertEquals("test-key", viewModel.uiState.value.apiKey)
        assertEquals("18790", viewModel.uiState.value.httpPort)
    }

    @Nested
    inner class FieldUpdates {

        @Test
        fun `onApiKeyChange updates apiKey`() {
            viewModel.onApiKeyChange("new-key")

            assertEquals("new-key", viewModel.uiState.value.apiKey)
        }

        @Test
        fun `onApiKeyChange clears error`() {
            viewModel.onApiKeyChange("key")

            assertNull(viewModel.uiState.value.error)
        }

        @Test
        fun `onHttpPortChange updates port`() {
            viewModel.onHttpPortChange("9090")

            assertEquals("9090", viewModel.uiState.value.httpPort)
        }

        @Test
        fun `onHttpPortChange rejects non-numeric`() {
            viewModel.onHttpPortChange("abc")

            assertEquals("18790", viewModel.uiState.value.httpPort)
        }

        @Test
        fun `onHttpPortChange accepts empty`() {
            viewModel.onHttpPortChange("")

            assertEquals("", viewModel.uiState.value.httpPort)
        }
    }

    @Nested
    inner class Save {

        @Test
        fun `save calls configApiClient and settingsStore`() = runTest {
            coEvery { configApiClient.saveConfig(any()) } returns SaveConfigResult(status = "ok")

            var completed = false
            viewModel.save { completed = true }
            advanceUntilIdle()

            assertTrue(completed)
            assertFalse(viewModel.uiState.value.saving)
            coVerify { configApiClient.saveConfig(any()) }
            coVerify { settingsStore.update(any()) }
        }

        @Test
        fun `save does not call API when hasErrors`() = runTest {
            // 99999 parses as Int but is out of 1..65535 range, so hasErrors=true
            viewModel.onHttpPortChange("99999")

            var completed = false
            viewModel.save { completed = true }
            advanceUntilIdle()

            assertFalse(completed)
            coVerify(exactly = 0) { configApiClient.saveConfig(any()) }
        }

        @Test
        fun `save sets error on exception`() = runTest {
            coEvery { configApiClient.saveConfig(any()) } throws RuntimeException("Network error")

            viewModel.save { }
            advanceUntilIdle()

            assertEquals("Network error", viewModel.uiState.value.error)
            assertFalse(viewModel.uiState.value.saving)
        }

        @Test
        fun `save with localOnly skips API call`() = runTest {
            val localOnlyHandle = SavedStateHandle(mapOf("localOnly" to true))
            val localVm = AppSettingsViewModel(localOnlyHandle, settingsStore, configApiClient)

            var completed = false
            localVm.save { completed = true }
            advanceUntilIdle()

            assertTrue(completed)
            coVerify(exactly = 0) { configApiClient.saveConfig(any()) }
            coVerify { settingsStore.update(any()) }
        }
    }

    private fun assertTrue(value: Boolean) {
        org.junit.jupiter.api.Assertions.assertTrue(value)
    }
}
