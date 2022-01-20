<template>
  <div>
    <div class="row flex q-pa-sm">
      <a
        v-if="this.$store.state.user.id == this.userID"
        ref="target_element"
        href="/upload"
        target="_blank"
        style="color: black"
        class="user-image-card add-link q-mr-sm q-mt-sm"
      >
        <img
          class="add-button"
          src="other/addBtn.png"
          ratio="1"
          spinner-color="orange"
        />
      </a>
      <user-image-card
        class="user-image-card q-mr-sm q-mt-sm"
        v-for="image in imageList"
        :inputImgUrl="image.md_url"
        :inputImgHashID="image.hash_id"
        :inputImgTitle="image.title"
        :key="image.key"
      >
      </user-image-card>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@import '../../css/width.scss';

.user-image-card {
  //顯示 2 張圖片
  @include xs-width {
    width: 47.5%;
  }

  //顯示 3 張圖片
  @include sm-width {
    width: 31.9%;
  }

  //顯示 4 張圖片
  @include md-width {
    width: 24%;
  }

  //顯示 5 張圖片
  @include lg-width {
    width: 19%;
  }

  //顯示 6 張圖片
  @include xl-width {
    width: 16.1%;
  }
}

.add-link {
  text-decoration: none;
}

.add-button {
  width: 100%;
  height: 100%;
}
</style>

<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useRoute } from 'vue-router';
import { api } from 'boot/axios';
import UserImageCard from 'src/components/User/ImageCard.vue';

export default defineComponent({
  name: 'PageUserImages',
  components: {
    'user-image-card': UserImageCard,
  },
  setup() {
    interface imageJson {
      key: number;
      hash_id: string;
      title: string;
      description: string;
      original_url: string;
      md_url?: string;
    }
    const route = useRoute();
    const userID = ref(parseInt(route.params.id as string));
    const imageList = ref([] as imageJson[]);

    async function getUserImages() {
      let path = '/user/' + userID.value.toString() + '/images';
      await api.get(path).then((res) => {
        let data = res.data as imageJson[];
        if (data) {
          imageList.value = data;
          for (let i = 0; i < imageList.value.length; i++) {
            if (!imageList.value[i].md_url) {
              imageList.value[i].md_url = imageList.value[i].original_url;
            }
            imageList.value[i].key = i;
          }
        }
      });
    }
    void getUserImages();

    return {
      userID,
      imageList,
    };
  },
});
</script>
