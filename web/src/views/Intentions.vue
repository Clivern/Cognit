<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="text-3xl font-semibold text-theme-text">{{ $t('intentions_page.title') }}</h1>
          <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
            {{ $t('dashboard.mock_badge') }}
          </span>
        </div>
        <p class="text-theme-textLight mt-2">{{ $t('intentions_page.subtitle') }}</p>
      </div>

      <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4">
          <h2 class="text-lg font-semibold text-theme-text">{{ $t('intentions_page.list_title') }}</h2>
          <p class="mt-1 text-sm text-theme-textLight">{{ $t('intentions_page.list_desc') }}</p>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('intentions_page.source') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('intentions_page.destination') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('intentions_page.skill') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('intentions_page.action') }}</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-theme-border">
              <tr v-for="rule in catalogIntentions" :key="rule.id" class="hover:bg-theme-hover">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium font-mono">
                  <router-link
                    v-if="agentExists(rule.source)"
                    :to="`/agents/${rule.source}`"
                    class="text-primary-800 hover:underline"
                  >
                    {{ rule.source }}
                  </router-link>
                  <span v-else class="text-theme-text">{{ rule.source }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium font-mono">
                  <router-link
                    v-if="agentExists(rule.destination)"
                    :to="`/agents/${rule.destination}`"
                    class="text-primary-800 hover:underline"
                  >
                    {{ rule.destination }}
                  </router-link>
                  <span v-else class="text-theme-text">{{ rule.destination }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ rule.skill }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="rule.action === 'allow' ? 'bg-emerald-50 text-emerald-800' : 'bg-red-50 text-red-800'"
                  >
                    {{ $t(`intentions_page.action_${rule.action}`) }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import AppNav from '@/components/AppNav.vue'
import { agentExists, catalogIntentions } from '@/mocks/catalog'
</script>
