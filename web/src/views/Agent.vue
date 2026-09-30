<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <router-link to="/agents" class="text-sm font-medium text-primary-800 hover:underline">
        {{ $t('agents_page.back_to_agents') }}
      </router-link>

      <div v-if="!agent" class="mt-8 rounded-xl border border-theme-border bg-white p-8 shadow-sm">
        <h1 class="text-2xl font-semibold text-theme-text">{{ $t('agents_page.not_found') }}</h1>
      </div>

      <template v-else>
        <div class="mt-4 mb-8">
          <div class="flex flex-wrap items-center gap-2">
            <h1 class="text-3xl font-semibold text-theme-text font-mono">{{ agent.name }}</h1>
            <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
              {{ $t('dashboard.mock_badge') }}
            </span>
            <span
              class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium"
              :class="healthClass(agent.health)"
            >
              <span class="h-1.5 w-1.5 rounded-full" :class="healthDotClass(agent.health)" aria-hidden="true" />
              {{ $t(`agents_page.health_${agent.health}`) }}
            </span>
          </div>
          <p class="text-theme-textLight mt-2">{{ $t(agent.descriptionKey) }}</p>
        </div>

        <div class="space-y-6">
          <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
            <div class="border-b border-theme-border px-6 py-4">
              <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.card') }}</h2>
              <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.card_desc') }}</p>
            </div>
            <dl class="divide-y divide-theme-border">
              <div class="grid gap-2 px-6 py-4 sm:grid-cols-[180px_1fr]">
                <dt class="text-sm text-theme-textLight">{{ $t('agents_page.version') }}</dt>
                <dd class="text-sm font-mono text-theme-text">{{ agent.version }}</dd>
              </div>
              <div class="grid gap-2 px-6 py-4 sm:grid-cols-[180px_1fr]">
                <dt class="text-sm text-theme-textLight">{{ $t('agents_page.skills') }}</dt>
                <dd class="flex flex-wrap gap-1.5">
                  <span
                    v-for="skill in agent.skills"
                    :key="skill.id"
                    class="rounded-full bg-theme-hover px-2 py-0.5 text-xs font-medium text-theme-text font-mono"
                  >
                    {{ skill.id }}
                  </span>
                </dd>
              </div>
              <div class="grid gap-2 px-6 py-4 sm:grid-cols-[180px_1fr]">
                <dt class="text-sm text-theme-textLight">{{ $t('agents_page.input_modes') }}</dt>
                <dd class="text-sm text-theme-text">{{ agent.inputModes.join(', ') }}</dd>
              </div>
              <div class="grid gap-2 px-6 py-4 sm:grid-cols-[180px_1fr]">
                <dt class="text-sm text-theme-textLight">{{ $t('agents_page.output_modes') }}</dt>
                <dd class="text-sm text-theme-text">{{ agent.outputModes.join(', ') }}</dd>
              </div>
              <div class="grid gap-2 px-6 py-4 sm:grid-cols-[180px_1fr]">
                <dt class="text-sm text-theme-textLight">{{ $t('agents_page.capabilities') }}</dt>
                <dd class="text-sm text-theme-text">{{ agent.capabilities.length ? agent.capabilities.join(', ') : '—' }}</dd>
              </div>
            </dl>
          </section>

          <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
            <div class="border-b border-theme-border px-6 py-4">
              <h2 class="text-lg font-semibold text-theme-text">{{ $t('agents_page.instances') }}</h2>
              <p class="mt-1 text-sm text-theme-textLight">{{ $t('agents_page.instances_desc') }}</p>
            </div>
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-theme-border">
                <thead class="bg-theme-hover">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.instance') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.address') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.datacenter') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('agents_page.health') }}</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-theme-border">
                  <tr v-for="instance in agent.instances" :key="instance.id">
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-theme-text font-mono">{{ instance.id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ instance.address }}:{{ instance.port }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text">{{ instance.datacenter }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                      <span
                        class="inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium"
                        :class="healthClass(instance.status)"
                      >
                        <span class="h-1.5 w-1.5 rounded-full" :class="healthDotClass(instance.status)" aria-hidden="true" />
                        {{ $t(`agents_page.health_${instance.status}`) }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import AppNav from '@/components/AppNav.vue'
import { findAgent } from '@/mocks/catalog'

const route = useRoute()
const agent = computed(() => findAgent(route.params.name))

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
