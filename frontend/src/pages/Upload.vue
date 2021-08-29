<template>
  <q-page class="flex flex-center">
    <loading-view :inputVisible="showLoading" />
    <q-dialog v-model="showMsg">
      <q-card>
        <q-card-section>
          <div class="text-h6 text-center">{{ msgTitle }}</div>
        </q-card-section>

        <q-card-section class="q-pt-none text-center">
          {{ msgContent }}
        </q-card-section>

        <q-card-actions align="center">
          <q-btn
            flat
            label="確認"
            color="primary"
            v-close-popup
            @click="redirectHomePage"
          />
        </q-card-actions>
      </q-card>
    </q-dialog>
    <q-uploader
      v-if="showPage"
      class="upload-block"
      url="/api/image"
      accept="image/jpeg, image/png, image/gif"
      label="圖片上傳"
      method="POST"
      field-name="image[]"
      multiple
      bordered
    />
  </q-page>
</template>

<style lang="sass" scoped>
.upload-block
  width: 300px
  height: 600px
</style>

<script>
import LoadingView from "../components/LoadingView.vue";

export default {
  name: "PageUpload",
  props: [],
  data() {
    return {
      userState: this.$store.state.user,
      showLoading: true,
      showMsg: false,
      msgTitle: "",
      msgContent: "",
      showPage: false
    };
  },
  components: {
    "loading-view": LoadingView
  },
  created() {
    (this.userState = this.$store.state.user), this.pageLoad();
  },
  watch: {
    "$store.state.user.dataLoaded": function() {
      this.pageLoad();
    }
  },
  methods: {
    pageLoad() {
      if (this.userState.dataLoaded) {
        this.showLoading = false;
        if (this.userState.id < 0) {
          // not login
          this.msgTitle = "您尚未登入";
          this.msgContent = "需要登入才能上傳圖片";
          this.showMsg = true;
        } else {
          this.showPage = true;
        }
      } else {
        this.showLoading = true;
      }
    },
    redirectHomePage() {
      this.$router.push("/");
    }
  }
};
</script>
