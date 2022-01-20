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
        <div class="text-center text-h5">您尋找的用戶不存在或是已刪除</div>
        <q-space style="height: 20px" />
        <div class="row justify-around not-found-redirect-link">
          <div class="text-h6 text-center" @click="push.homePage()">
            返回首頁
          </div>
          <div class="text-h6 text-center" @click="push.previousPage()">
            上一頁
          </div>
        </div>
      </div>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.user-card-row {
  width: 100%;
}

.tab-row {
  width: 100%;
  height: 50px;
}

.not-found-colunm {
  width: 100%;
}

.not-found-msg {
  margin-top: 24px;
  margin-bottom: 24px;
  width: 350px;
}

.not-found-redirect-link {
  color: #337f15;
}
</style>

<script lang="ts">
import { ref, defineComponent } from 'vue';
import { useRoute } from 'vue-router';
import axios from 'axios';
import { api } from 'boot/axios';
import { Push } from 'src/lib/router/pushPage';
import LoadingView from '../components/LoadingView.vue';

export default defineComponent({
  name: 'PageUser',
  components: {
    'loading-view': LoadingView,
  },
  setup() {
    const route = useRoute();
    const push = new Push();

    const userID = ref(route.params.id);
    const userShowName = ref('');
    const userAvatar = ref('');
    const userIntroduction = ref('');
    const tabRef = ref('images');
    const showLoading = ref(true);
    const userDataLoaded = ref(false);
    const showNotFound = ref(false);

    function initStatus() {
      userDataLoaded.value = false;
      showLoading.value = true;
      showNotFound.value = false;
    }
    function notFoundPage() {
      showLoading.value = false;
      showNotFound.value = true;
    }
    function normalPage() {
      showLoading.value = false;
      userDataLoaded.value = true;
    }
    initStatus();

    interface userDataJSON {
      show_name: string;
      avatar: string;
      introduction: string;
    }
    async function getUserInfo(id?: number) {
      if (id) {
        var path = '/user/' + id.toString();
      } else {
        var path = '/user/' + userID.value.toString();
      }
      await api
        .get(path)
        .then((res) => {
          var data = res.data as userDataJSON;
          userShowName.value = data['show_name'];
          userAvatar.value = data['avatar'];
          userIntroduction.value = data['introduction'];
          // 載入正常頁面
          normalPage();
        })
        .catch((error) => {
          if (axios.isAxiosError(error)) {
            if (error.response) {
              switch (error.response.status) {
                // 400 -> ID 不是數字
                // 404 -> 找不到用戶
                case 400:
                case 404:
                  notFoundPage();
                  break;
              }
            }
          }
        });
    }
    void getUserInfo();

    return {
      userID,
      userShowName,
      userAvatar,
      userIntroduction,
      tabRef,
      showLoading,
      initStatus,
      notFoundPage,
      normalPage,
      getUserInfo,
      userDataLoaded,
      showNotFound,
      push,
    };
  },
  beforeRouteUpdate(to, _, next) {
    // 重置資料載入狀態，以顯示 loading 效果
    this.initStatus();
    const idStr = to.params.id as string;
    void this.getUserInfo(Number.parseInt(idStr));
    next();
  },
  methods: {},
});
</script>
