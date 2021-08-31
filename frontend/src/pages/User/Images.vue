<template>
  <div>
    <div class="row flex q-pa-sm">
      <a
        v-if="this.$store.state.user.id == this.userID"
        ref="target_element"
        href="/upload"
        target="_blank"
        style="color: black;"
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
@import "../../css/width.scss";

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

<script>
import UserImageCard from "../../components/User/ImageCard.vue";

export default {
  name: "PageUserImages",
  components: {
    "user-image-card": UserImageCard
  },
  props: [],
  data() {
    return {
      imageHeight: 0,
      userID: this.$route.params.id,
      imageList: []
    };
  },
  created() {
    this.getUserImages();
  },
  methods: {
    async getUserImages() {
      var path = "/api/user/" + this.userID + "/images";
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.imageList = data;
          for (var i = 0; i < this.imageList.length; i++) {
            if (this.imageList[i].md_url == null) {
              this.imageList[i].md_url = this.imageList[i].original_url;
            }
            this.imageList[i].key = i;
          }
        })
        .catch(error => {
          if (error.request) {
          } else if (error.response) {
            switch (error.response.status) {
              //Internal Server Error
              case 500:
                break;
              //Unauthorized
              case 401:
                break;
              //Not Found
              case 404:
                break;
            }
          }
        });
    }
  }
};
</script>
