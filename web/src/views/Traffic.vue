<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="text-3xl font-semibold text-theme-text">{{ $t('traffic_page.title') }}</h1>
          <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
            {{ $t('dashboard.mock_badge') }}
          </span>
        </div>
        <p class="text-theme-textLight mt-2">{{ $t('traffic_page.subtitle') }}</p>
      </div>

      <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4">
          <h2 class="text-lg font-semibold text-theme-text">{{ $t('traffic_page.list_title') }}</h2>
          <p class="mt-1 text-sm text-theme-textLight">{{ $t('traffic_page.list_desc') }}</p>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.when') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.source') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.destination') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.skill') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.verb') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.decision') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.result') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.instance') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.latency') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('traffic_page.task') }}</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-theme-border">
              <tr v-for="call in catalogTraffic" :key="call.id" class="hover:bg-theme-hover">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ $t(call.timeKey) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium font-mono">
                  <router-link
                    v-if="agentExists(call.source)"
                    :to="`/agents/${call.source}`"
                    class="text-primary-800 hover:underline"
                  >
                    {{ call.source }}
                  </router-link>
                  <span v-else class="text-theme-text">{{ call.source }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium font-mono">
                  <router-link
                    v-if="agentExists(call.destination)"
                    :to="`/agents/${call.destination}`"
                    class="text-primary-800 hover:underline"
                  >
                    {{ call.destination }}
                  </router-link>
                  <span v-else class="text-theme-text">{{ call.destination }}</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ call.skill }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ call.verb }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="call.decision === 'allow' ? 'bg-emerald-50 text-emerald-800' : 'bg-red-50 text-red-800'"
                  >
                    {{ $t(`traffic_page.decision_${call.decision}`) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    class="inline-flex rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="resultClass(call.result)"
                  >
                    {{ $t(`traffic_page.result_${call.result}`) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ call.instance || '—' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text">{{ call.latency || '—' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ call.taskId || '—' }}</td>
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
import { agentExists, catalogTraffic } from '@/mocks/catalog'

function resultClass(result) {
  return {
    completed: 'bg-emerald-50 text-emerald-800',
    failed: 'bg-red-50 text-red-800',
    denied: 'bg-red-50 text-red-800',
    timeout: 'bg-amber-50 text-amber-800',
    no_instance: 'bg-amber-50 text-amber-800',
  }[result] || 'bg-theme-hover text-theme-textLight'
}
</script>
