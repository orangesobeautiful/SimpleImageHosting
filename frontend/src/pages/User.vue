<template>
  <q-page>
    <!-- loading page -->
    <loading-view :inputVisible="showLoading" />
    <!--normal page-->
    <div v-if="userDataLoaded">
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
    </div>
    <!-- error page -->
    <div v-if="showNotFound" class="column not-found-colunm items-center">
      <div class="not-found-msg">
        <div class="text-center text-h5">
          您尋找的用戶不存在或是已刪除
        </div>
        <q-space style="height: 20px" />
        <div class="row justify-around not-found-redirect-link">
          <div class="text-h6 text-center" @click="goHomePage">
            返回首頁
          </div>
          <div class="text-h6 text-center" @click="goPreviousPage">
            上一頁
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<style lang="sass" scoped>
.user-card-row
  width: 100%
.tab-row
  width: 100%
  height: 50px
.not-found-colunm
  width: 100%
.not-found-msg
  margin-top: 24px
  margin-buttom: 24px
  width: 350px
.not-found-redirect-link
  color: #337f15
</style>

<script>
import LoadingView from "../components/LoadingView.vue";

export default {
  name: "PageUser",
  data() {
    return {
      userID: this.$route.params.id,
      userShowName: "",
      userAvatar: "",
      userIntroduction: "",
      tabRef: "images",
      showLoading: true,
      userDataLoaded: false,
      showNotFound: false
    };
  },
  components: {
    "loading-view": LoadingView
  },
  created() {
    this.initStatus();
    this.getUserInfo();
  },
  beforeRouteUpdate(to, _, next) {
    // 重置資料載入狀態，以顯示 loading 效果
    this.initStatus();
    this.getUserInfo(to.params.id);
    next();
  },
  methods: {
    initStatus() {
      this.userDataLoaded = false;
      this.showLoading = true;
      this.showNotFound = false;
    },
    notFoundPage() {
      this.showLoading = false;
      this.showNotFound = true;
    },
    normalPage() {
      this.showLoading = false;
      this.userDataLoaded = true;
    },
    async getUserInfo(id = null) {
      if (id == null) {
        var path = "/api/user/" + this.userID;
      } else {
        var path = "/api/user/" + id;
      }

      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.userShowName = data["show_name"];
          this.userAvatar = data["avatar"];
          this.userIntroduction = data["introduction"];

          this.normalPage();
        })
        .catch(error => {
          if (error.request) {
            switch (error.request.status) {
              //Internal Server Error
              case 500:
                break;
              //Unauthorized
              case 401:
                break;
              //Bad Request
              //Not Found
              case 400:
              case 404:
                this.notFoundPage();
                break;
            }
          } else if (error.response) {
          }
        });
    },
    goHomePage() {
      this.$router.push("/");
    },
    goPreviousPage() {
      this.$router.go(-1);
    }
  }
};
</script>
