<template>
  <div class="min-h-screen bg-theme-bg">
    <AppNav />

    <main class="w-full py-8 px-4 sm:px-6 lg:px-8">
      <div class="mb-8">
        <div class="flex flex-wrap items-center gap-2">
          <h1 class="text-3xl font-semibold text-theme-text">{{ $t('kv_page.title') }}</h1>
          <span class="rounded-full border border-amber-200 bg-amber-50 px-2.5 py-0.5 text-xs font-medium text-amber-800">
            {{ $t('dashboard.mock_badge') }}
          </span>
        </div>
        <p class="text-theme-textLight mt-2">{{ $t('kv_page.subtitle') }}</p>
      </div>

      <section class="rounded-xl border border-theme-border bg-white shadow-sm overflow-hidden">
        <div class="border-b border-theme-border px-6 py-4">
          <h2 class="text-lg font-semibold text-theme-text">{{ $t('kv_page.list_title') }}</h2>
          <p class="mt-1 text-sm text-theme-textLight">{{ $t('kv_page.list_desc') }}</p>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-theme-border">
            <thead class="bg-theme-hover">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('kv_page.key') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('kv_page.value') }}</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-theme-textLight uppercase tracking-wider">{{ $t('kv_page.updated') }}</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-theme-border">
              <tr v-for="item in catalogKV" :key="item.key" class="hover:bg-theme-hover">
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-theme-text font-mono">{{ item.key }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-text font-mono">{{ item.value }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-theme-textLight">{{ $t(item.updatedKey) }}</td>
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
import { catalogKV } from '@/mocks/catalog'
</script>
