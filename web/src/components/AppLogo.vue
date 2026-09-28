<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="5 5 501 597"
    fill="currentColor"
    :role="alt ? 'img' : 'presentation'"
    :aria-label="alt || undefined"
    :aria-hidden="alt ? undefined : true"
  >
    <path d="M180.5 6 188 5.5 194.5 12 332 447.5 401.5 282 407 274.5 498 273.5 505.5 280 505.5 327 502 331.5 498 333.5 447 333.5 444.5 335 336.5 592 333.5 598 329 601.5 322 601.5 315.5 595 178 161.5 108.5 328 104 332.5 12 333.5 5.5 327 5.5 281 9 275.5 14 273.5 65 273.5 66.5 272 174.5 12Z" />
  </svg>
</template>

<script setup>
defineProps({
  alt: {
    type: String,
    default: ''
  }
})
</script>
