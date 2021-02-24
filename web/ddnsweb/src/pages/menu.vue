<template>
  <div id="app">
    <nav-menu 
    :themes="themes" 
    :menus="menus" 
    :logo="logo" 
    :systemName="systemName" 
    ref="NewTap"> </nav-menu>
  </div>
</template>
<script>
import navMenu from 'nav-menu'; // 引入
export default {
  name: 'app',
  data() {
    return {
      logo: "http://bj.images.cdqykj.cn/sso/4cd16a1b5b36ac980713ce25c3c7ab0f.png",
      themes: "bg-info|bg-info|bg-dark dark-danger", //顶部左侧背景颜色,顶部右侧背景颜色,右边菜单背景颜色
      systemName: "千行网址系统",  //系统名称
      userinfo: { name: '', role: "管理员" },
      indexUrl: "/home",
      items: [],
      menus:[{}],
    }
  },
  components: { //注册插件
    navMenu
  },
  created() {
    //  this.getMenu();
    // this.getSystemInfo();
  },
  mounted() {
    this.setDocmentTitle();
    // this.userinfo = JSON.parse(localStorage.getItem("userinfo"));
    this.$refs.NewTap.open("用户管理", this.indexUrl);
  },
  methods: {
    setDocmentTitle() {
      document.title = this.systemName;
    }

  }
}
</script>

<style>
.nav-tabs {
  display: none;
}
.app-aside {
  display: none;
  width: 0;
}
.app-content,
.app-footer {
  margin-left: 0px !important;
  text-align: center;
}
.app-content-body{
  padding-bottom: 0px !important;
}
.navbar-nav {
  display: none;
}
</style>
