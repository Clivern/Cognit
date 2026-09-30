<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="text-3xl font-semibold text-theme-text">{{ $t('agents_page.title') }}</h1>
          <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
            {{ $t('dashboard.mock_badge') }}
          </span>
        </div>
        <p class="text-theme-textLight mt-2">{{ $t('agents_page.subtitle') }}</p>
      </div>

      <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4">
          <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.list_title') }}</h2>
          <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.list_desc') }}</p>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.name') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.version') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.skills') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.instances') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.health') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.updated') }}</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-theme-border">
              <tr v-for="agent in agents" :key="agent.id" class="hover:bg-theme-hover">
                <td class="px-6 py-4">
                  <p class="text-sm font-medium text-theme-text font-mono">{{ agent.name }}</p>
                  <p class="mt-0.5 text-sm text-theme-textLight">{{ agent.description }}</p>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ agent.version }}</td>
                <td class="px-6 py-4">
                  <div class="flex flex-wrap gap-1.5">
                    <span
                      v-for="skill in agent.skills"
                      :key="skill"
                      class="rounded-full bg-theme-hover px-2 py-0.5 text-xs font-medium text-theme-text font-mono"
                    >
                      {{ skill }}
                    </span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text">
                  {{ $t('agents_page.live_count', { count: agent.live }) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="healthClass(agent.health)"
                  >
                    <span class="h-1.5 w-1.5 rounded-full" :class="healthDotClass(agent.health)" aria-hidden="true" />
                    {{ $t(`agents_page.health_${agent.health}`) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ agent.updated }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import AppNav from '@/components/AppNav.vue'

const { t } = useI18n()

const agents = computed(() => [
  {
    id: 'support-assistant',
    name: 'support-assistant',
    description: t('agents_page.mock_support_desc'),
    version: '1.0.0',
    skills: ['support.answer'],
    live: 2,
    health: 'passing',
    updated: t('agents_page.mock_updated_2m'),
  },
  {
    id: 'user-directory',
    name: 'user-directory',
    description: t('agents_page.mock_directory_desc'),
    version: '1.2.0',
    skills: ['user.profile.get'],
    live: 1,
    health: 'passing',
    updated: t('agents_page.mock_updated_18m'),
  },
  {
    id: 'invoice-extractor',
    name: 'invoice-extractor',
    description: t('agents_page.mock_invoice_desc'),
    version: '0.9.1',
    skills: ['invoice.extract'],
    live: 1,
    health: 'warning',
    updated: t('agents_page.mock_updated_1h'),
  },
])

function healthClass(health) {
  return {
    passing: 'bg-emerald-50 text-emerald-800',
    warning: 'bg-amber-50 text-amber-800',
    critical: 'bg-red-50 text-red-800',
  }[health] || 'bg-theme-hover text-theme-textLight'
}

function healthDotClass(health) {
  return {
    passing: 'bg-emerald-500',
    warning: 'bg-amber-500',
    critical: 'bg-red-500',
  }[health] || 'bg-theme-textLight'
}
</script>
