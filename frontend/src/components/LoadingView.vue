<template>
  <q-inner-loading :showing="inputVisible">
    <q-spinner-ios size="100px" color="brown" />
  </q-inner-loading>
</template>

<script>
export default {
  name: "LoadingView",
  props: ["inputVisible"],
  data() {
    return {};
  },
  created() {},
  methods: {}
};
</script>
