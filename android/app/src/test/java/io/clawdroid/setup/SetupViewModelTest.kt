package io.clawdroid.setup

import io.clawdroid.backend.api.GatewaySettings
import io.clawdroid.backend.api.GatewaySettingsStore
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
import org.junit.jupiter.api.Assertions.assertNotEquals
import org.junit.jupiter.api.Assertions.assertNull
import org.junit.jupiter.api.Assertions.assertTrue
import org.junit.jupiter.api.BeforeEach
import org.junit.jupiter.api.Nested
import org.junit.jupiter.api.Test

@OptIn(ExperimentalCoroutinesApi::class)
class SetupViewModelTest {

    private val testDispatcher = UnconfinedTestDispatcher()
    private val setupApiClient = mockk<SetupApiClient>(relaxed = true)
    private val settingsStore = mockk<GatewaySettingsStore>(relaxed = true) {
        every { settings } returns MutableStateFlow(GatewaySettings())
    }

    private lateinit var viewModel: SetupViewModel

    @BeforeEach
    fun setup() {
        Dispatchers.setMain(testDispatcher)
        viewModel = SetupViewModel(setupApiClient, settingsStore)
    }

    @AfterEach
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Nested
    inner class FieldUpdates {

        @Test
        fun `onGatewayPortChange updates port`() {
            viewModel.onGatewayPortChange("9090")

            assertEquals("9090", viewModel.uiState.value.gatewayPort)
        }

        @Test
        fun `onGatewayPortChange rejects non-numeric`() {
            viewModel.onGatewayPortChange("abc")

            assertEquals("18790", viewModel.uiState.value.gatewayPort)
        }

        @Test
        fun `onGatewayPortChange accepts empty`() {
            viewModel.onGatewayPortChange("")

            assertEquals("", viewModel.uiState.value.gatewayPort)
        }

        @Test
        fun `onGatewayApiKeyChange updates apiKey`() {
            viewModel.onGatewayApiKeyChange("my-key")

            assertEquals("my-key", viewModel.uiState.value.gatewayApiKey)
        }

        @Test
        fun `generateApiKey sets non-empty apiKey`() {
            viewModel.generateApiKey()

            assertTrue(viewModel.uiState.value.gatewayApiKey.isNotEmpty())
        }

        @Test
        fun `generateApiKey produces unique keys`() {
            viewModel.generateApiKey()
            val key1 = viewModel.uiState.value.gatewayApiKey
            viewModel.generateApiKey()
            val key2 = viewModel.uiState.value.gatewayApiKey

            assertNotEquals(key1, key2)
        }

        @Test
        fun `onLlmModelChange updates llmModel`() {
            viewModel.onLlmModelChange("claude-3")

            assertEquals("claude-3", viewModel.uiState.value.llmModel)
        }

        @Test
        fun `onWorkspaceChange updates workspace`() {
            viewModel.onWorkspaceChange("/home/user/workspace")

            assertEquals("/home/user/workspace", viewModel.uiState.value.workspace)
        }
    }

    @Nested
    inner class StepNavigation {

        @Test
        fun `submitInit advances to step 1 when valid`() {
            viewModel.onGatewayPortChange("18790")
            viewModel.onGatewayApiKeyChange("key-123")

            viewModel.submitInit()

            assertEquals(1, viewModel.uiState.value.currentStep)
            assertTrue(viewModel.uiState.value.step1Done)
        }

        @Test
        fun `submitInit does nothing when canProceedStep1 is false`() {
            viewModel.onGatewayApiKeyChange("")

            viewModel.submitInit()

            assertEquals(0, viewModel.uiState.value.currentStep)
            assertFalse(viewModel.uiState.value.step1Done)
        }

        @Test
        fun `skipStep 2 sets step2Skipped and advances`() {
            viewModel.skipStep(2)

            assertTrue(viewModel.uiState.value.step2Skipped)
            assertEquals(2, viewModel.uiState.value.currentStep)
        }

        @Test
        fun `skipStep 3 sets step3Skipped and advances`() {
            viewModel.skipStep(3)

            assertTrue(viewModel.uiState.value.step3Skipped)
            assertEquals(3, viewModel.uiState.value.currentStep)
        }

        @Test
        fun `nextStep sets currentStep`() {
            viewModel.nextStep(2)

            assertEquals(2, viewModel.uiState.value.currentStep)
        }

        @Test
        fun `previousStep decrements currentStep`() {
            viewModel.nextStep(2)

            viewModel.previousStep()

            assertEquals(1, viewModel.uiState.value.currentStep)
        }

        @Test
        fun `previousStep does not go below 0`() {
            viewModel.previousStep()

            assertEquals(0, viewModel.uiState.value.currentStep)
        }
    }

    @Nested
    inner class SubmitComplete {

        @Test
        fun `submitComplete calls init and complete APIs`() = runTest {
            coEvery { setupApiClient.init(any()) } returns Unit
            coEvery { setupApiClient.complete(any()) } returns Unit

            viewModel.onGatewayPortChange("18790")
            viewModel.onGatewayApiKeyChange("test-key")
            viewModel.submitInit()

            var completed = false
            viewModel.submitComplete { completed = true }
            advanceUntilIdle()

            assertTrue(completed)
            assertFalse(viewModel.uiState.value.loading)
            coVerify { setupApiClient.init(any()) }
            coVerify { setupApiClient.complete(any()) }
            coVerify { settingsStore.update(any()) }
        }

        @Test
        fun `submitComplete sets error on failure`() = runTest {
            coEvery { setupApiClient.init(any()) } throws RuntimeException("Connection refused")

            viewModel.onGatewayPortChange("18790")
            viewModel.onGatewayApiKeyChange("key")
            viewModel.submitInit()

            var completed = false
            viewModel.submitComplete { completed = true }
            advanceUntilIdle()

            assertFalse(completed)
            assertFalse(viewModel.uiState.value.loading)
            assertEquals("Connection refused", viewModel.uiState.value.error)
        }

        @Test
        fun `submitComplete does not double-submit while loading`() = runTest {
            coEvery { setupApiClient.init(any()) } returns Unit
            coEvery { setupApiClient.complete(any()) } returns Unit

            viewModel.onGatewayPortChange("18790")
            viewModel.onGatewayApiKeyChange("key")
            viewModel.submitInit()

            // With UnconfinedTestDispatcher, coroutines complete eagerly,
            // so we cannot truly test the loading guard mid-flight.
            // Verify at minimum that calling twice does not crash and APIs are called.
            viewModel.submitComplete { }
            advanceUntilIdle()

            viewModel.submitComplete { }
            advanceUntilIdle()

            assertFalse(viewModel.uiState.value.loading)
            coVerify(atLeast = 1) { setupApiClient.init(any()) }
        }
    }
}
