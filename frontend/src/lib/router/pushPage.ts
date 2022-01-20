import { useRouter, Router } from 'vue-router';

export class Push {
  private router: Router;

  constructor() {
    this.router = useRouter();
  }

  // 跳轉到指定頁面
  public link(url: string) {
    void this.router.push(url);
  }

  // 跳轉至 登入頁面
  public homePage() {
    void this.router.push('/');
  }

  // 跳轉至 登入頁面
  public signinPage() {
    void this.router.push('/signin');
  }

  // 跳轉至 註冊頁面
  public registerPage() {
    void this.router.push('/register');
  }

  // 跳轉至 上傳頁面
  public uploadPage() {
    void this.router.push('/upload');
  }

  // 跳轉至 網頁控制台頁面
  public dashboardPage() {
    void this.router.push('/dashboard/settings');
  }

  // 跳轉至 使用者圖片頁面
  public userImagesPage(userID: number) {
    void this.router.push('/user/' + userID.toString() + '/images');
  }

  // 重新整理
  public reload() {
    this.router.go(0);
  }

  // 上一頁
  public previousPage() {
    this.router.go(-1);
  }
}
