<template>
  <div
    class="image-button"
    @click="push.link(viewURL)"
    @click.middle="newBlankImagePage"
  >
    <q-img :src="imgURL" :ratio="1" spinner-color="orange">
      <div
        v-if="imgTitle"
        class="absolute-bottom-right text-subtitle1"
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

<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useRouter } from 'vue-router';
import { Push } from 'src/lib/router/pushPage';

export default defineComponent({
  name: 'imgCard',
  props: {
    inputImgUrl: String,
    inputImgHashID: String,
    inputImgTitle: String,
  },
  setup(props) {
    const router = useRouter();
    const push = new Push();

    const imgURL = ref(props.inputImgUrl);
    const hashID = ref(props.inputImgHashID);
    const imgTitle = ref(props.inputImgTitle);
    const viewURL = ref('');
    if (hashID.value) {
      viewURL.value = '/image/' + hashID.value;
    }

    function newBlankImagePage() {
      let route = router.resolve(viewURL.value);
      window.open(route.href);
    }

    return {
      push,
      imgURL,
      hashID,
      imgTitle,
      viewURL,
      newBlankImagePage,
    };
  },
  methods: {},
});
</script>
