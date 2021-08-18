<template>
  <q-page>
    <div class="row background-container">
      <q-img class="image-block" src="background/homepage-background.jpg">
        <div class="absolute-full flex flex-center">
          <div class="column items-center">
            <div class="text-h3 text-weight-bolder title">這是個示範用網站</div>
            <div class="space" />
            <div class="text-h5 description">
              此網站僅供示範用，您上傳的圖片、註冊的帳號隨時有可能被刪除，請不要用來儲存重要資料。
            </div>
          </div>
        </div>
      </q-img>
    </div>
  </q-page>
</template>

<style lang="scss" scoped>
.space {
  height: 20px;
}

.title {
  min-width: 384px;
}

.description {
  min-width: 395px;
}

.background-container {
  width: 100%;
}
.image-block {
  height: calc(100vh - 50px);
  width: 100vw;
}
</style>

<script>
export default {
  name: "PageIndex"
};
</script>
