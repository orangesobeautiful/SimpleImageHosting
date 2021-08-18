<template>
  <q-page>
    <div class="user-card-row column">
      <div class="row q-pa-md q-gutter-sm items-center">
        <q-avatar size="100px">
          <img :src="userAvatar" />
        </q-avatar>
        <div class="text-h3">{{ userShowName }}</div>
      </div>
      <q-tabs v-model="tabRef" inline-label class="text-black">
        <q-route-tab to="images" name="images" label="圖片" />
        <q-route-tab to="about" name="about" label="關於" />
      </q-tabs>
      <q-separator color="black" inset />
    </div>

    <router-view />
  </q-page>
</template>

<style lang="sass" scoped>
.user-card-row
  width: 100%
.tab-row
  width: 100%
  height: 50px
</style>

<script>
export default {
  name: "PageUser",
  data() {
    return {
      userID: this.$route.params.id,
      userShowName: "",
      userAvatar: "",
      userIntroduction: "",
      tabRef: "images"
    };
  },
  created() {
    this.getUserInfo();
  },
  methods: {
    async getUserInfo() {
      var path = "/api/user/" + this.userID;
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.userShowName = data["show_name"];
          this.userAvatar = data["avatar"];
          this.userIntroduction = data["introduction"];
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
