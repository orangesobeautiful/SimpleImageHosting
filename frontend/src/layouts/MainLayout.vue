<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated class="toolbar">
      <div>
        <q-toolbar class="text-black">
          <q-toolbar-title>
            <q-btn flat @click="homePage">
              Simple Image Hosting
            </q-btn>
          </q-toolbar-title>

          <div class="row" v-if="loginInfoChecked">
            <div v-if="!hasLogin">
              <q-btn @click="registerPage">Register</q-btn>
              <q-btn @click="signinPage">Signin</q-btn>
            </div>
            <div v-if="hasLogin" class="row items-center">
              <div class="row">{{ showName }}</div>
              <q-avatar class="q-ml-sm">
                <img :src="avatar" />
                <q-menu :style="{ backgroundColor: '#eee', color: 'blue' }">
                  <q-list style="min-width: 100px">
                    <q-item clickable @click="myImagesPage">
                      <q-item-section>我的圖片</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item clickable @click="uploadPage">
                      <q-item-section>上傳</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item v-if="grade == 1" clickable @click="dashboardPage">
                      <q-item-section>控制台</q-item-section>
                    </q-item>
                    <q-separator />
                    <q-item clickable @click="logout">
                      <q-item-section>登出</q-item-section>
                    </q-item>
                  </q-list>
                </q-menu>
              </q-avatar>
            </div>
          </div>
        </q-toolbar>
      </div>
    </q-header>

    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<style lang="sass" scoped>
.toolbar
  background-color: rgba(251, 115, 4, 0.79)
</style>

<script>
export default {
  name: "MainLayout",
  data() {
    return {
      leftDrawerOpen: false,
      loginInfoChecked: false,
      hasLogin: false,
      userID: 0,
      showName: "",
      avatar: "",
      grade: -1
    };
  },
  async created() {
    await this.getUserInfo();
    this.loginInfoChecked = true;
  },
  methods: {
    async getUserInfo() {
      var path = "/api/me";
      await this.$axios
        .get(path)
        .then(res => {
          var data = res.data;
          this.hasLogin = true;
          this.userID = data["id"];
          this.showName = data["show_name"];
          this.avatar = data["avatar"];
          this.grade = data["grade"];

          this.$store.commit("user/setID", this.userID);
          this.$store.commit("user/setGrade", this.grade);
          this.$store.commit("user/setDataLoaded", true);
        })
        .catch(error => {
          this.$store.commit("user/setDataLoaded", true);
          if (error.request) {
          } else if (error.response) {
            switch (error.response.status) {
              //Internal Server Error
              case 500:
                break;
              //Unauthorized
              case 401:
                this.hasLogin = false;
                break;
            }
          }
        });
    },
    async logout() {
      var path = "/api/logout";
      await this.$axios
        .post(path)
        .then(() => {
          this.$store.commit("user/setID", -1);
          this.$store.commit("user/setGrade", -1);
          this.reload();
        })
        .catch(() => {});
    },
    signinPage() {
      this.$router.push("/signin");
    },
    registerPage() {
      this.$router.push("/register");
    },
    myImagesPage() {
      this.$router.push("/user/" + this.userID.toString() + "/images");
    },
    uploadPage() {
      this.$router.push("/upload");
    },
    dashboardPage() {
      this.$router.push("/dashboard/settings");
    },
    homePage() {
      this.$router.push("/");
    },
    reload() {
      this.$router.go(0);
    }
  }
};
</script>
