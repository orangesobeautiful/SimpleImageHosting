<template>
  <div
    class="image-button"
    @click="goImagePage"
    @click.middle="newBlankImagePage"
  >
    <q-img :src="imgURL" :ratio="1" spinner-color="orange">
      <div
        v-if="imgTitle"
        class="absolute-bottom-right text-subtitle1 "
        style="border-top-left-radius: 3px"
      >
        {{ imgTitle }}
      </div>
    </q-img>
  </div>
</template>

<style lang="scss" scoped>
.image-button {
  cursor: pointer;
}
</style>

<script>
export default {
  name: "imgCard",
  props: ["inputImgUrl", "inputImgHashID", "inputImgTitle"],
  data() {
    return {
      imgURL: this.inputImgUrl,
      hashID: this.inputImgHashID,
      imgTitle: this.inputImgTitle,
      viewUrl: ""
    };
  },
  created() {
    this.viewUrl = "/image/" + this.hashID;
  },
  methods: {
    goImagePage() {
      this.$router.push(this.viewUrl);
    },
    newBlankImagePage() {
      let route = this.$router.resolve({ path: this.viewUrl });
      window.open(route.href);
    }
  }
};
</script>
