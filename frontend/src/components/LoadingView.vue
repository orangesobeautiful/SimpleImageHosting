<template>
  <q-inner-loading :showing="inputVisible">
    <q-spinner-ios size="100px" color="brown" />
  </q-inner-loading>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
export default defineComponent({
  name: 'LoadingView',
  props: ['inputVisible'],
});
</script>
