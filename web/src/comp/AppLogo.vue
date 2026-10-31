<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    :viewBox="variant === 'mask' ? '96 96 320 320' : '0 0 512 512'"
    fill="currentColor"
    :role="alt ? 'img' : 'presentation'"
    :aria-label="alt || undefined"
    :aria-hidden="alt ? undefined : true"
  >
    <template v-if="variant === 'mask'">
      <path
        d="M256 256L378.6 358.9M256 256L378.6 153.1"
        fill="none"
        stroke="currentColor"
        stroke-width="18"
        stroke-linecap="round"
      />
      <circle cx="256" cy="256" r="46" />
      <circle cx="378.6" cy="358.9" r="30" />
      <circle cx="378.6" cy="153.1" r="30" />
    </template>
    <template v-else-if="variant === 'tile'">
      <rect width="512" height="512" rx="116" />
      <g transform="translate(256 256) scale(.78) translate(-256 -256)">
        <g fill="none" stroke="var(--theme-bg)" stroke-linecap="round">
          <path d="M378.6 358.9A160 160 0 1 1 378.6 153.1" stroke-width="56" />
          <path d="M256 256L378.6 358.9M256 256L378.6 153.1" stroke-width="18" />
        </g>
        <g fill="var(--theme-bg)">
          <circle cx="256" cy="256" r="46" />
          <circle cx="378.6" cy="358.9" r="30" />
          <circle cx="378.6" cy="153.1" r="30" />
        </g>
      </g>
    </template>
    <template v-else>
      <g fill="none" stroke="currentColor" stroke-linecap="round">
        <path d="M378.6 358.9A160 160 0 1 1 378.6 153.1" stroke-width="56" />
        <path d="M256 256L378.6 358.9M256 256L378.6 153.1" stroke-width="18" />
      </g>
      <circle cx="256" cy="256" r="46" />
      <circle cx="378.6" cy="358.9" r="30" />
      <circle cx="378.6" cy="153.1" r="30" />
    </template>
  </svg>
</template>

<script setup>
defineProps({
  alt: {
    type: String,
    default: '',
  },
  variant: {
    type: String,
    default: 'current',
  },
})
</script>
