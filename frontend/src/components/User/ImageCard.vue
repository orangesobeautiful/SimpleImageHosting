<template>
  <a :href="viewUrl"
    ><q-img :src="imgURL" :ratio="1" spinner-color="orange"></q-img
  ></a>
</template>

<script>
export default {
  name: "imgCard",
  props: ["inputImgUrl", "inputImgHashID"],
  data() {
    return {
      imgURL: this.inputImgUrl,
      hashID: this.inputImgHashID,
      viewUrl: ""
    };
  },
  created() {
    this.viewUrl = "/image/" + this.hashID;
  },
  methods: {}
};
</script>
