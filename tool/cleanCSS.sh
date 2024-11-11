go install github.com/daaku/cssdalek@latest
cssdalek \
  --css '../public/assets/frontend/css/frontend.min.css'\
  --word '../public/assets/frontend/js/bootstrap/bootstrap.js'\
  --word '../template/frontend/default/*.html'\
  --word '../template/frontend/default/*/*.html'\
  --word '../template/frontend/default/*/*/*.html'\
  --word '../template/frontend/default/*/*/*/*.html'\
  --word '../public/assets/frontend/js/*.js'\
  --word '../public/assets/frontend/js/superfish/superfish.js'\
  --word '../public/assets/frontend/js/page/*.js'\
  --word '../public/assets/frontend/js/page/*/*.js'\
  --include-class "sf-.*"\
  --include-selector ".sf-menu.nav-info a"\
  --include-selector ".sf-menu.nav-info ul"\
  --include-selector "menu.nav-info ul li a"\
  --include-selector ".sf-menu.nav-info .active"\
  --include-id 'captchaImage' > ../public/assets/frontend/css/frontend.lite.min.css

cssdalek \
  --css '../public/assets/frontend/blogdefault/css/vendors/bootstrap.min.css'\
  --word '../public/assets/frontend/js/bootstrap/bootstrap.js'\
  --word '../template/frontend/blogdefault/*.html'\
  --word '../template/frontend/blogdefault/*/*.html'\
  --word '../template/frontend/blogdefault/*/*/*.html'\
  --word '../public/assets/frontend/js/*.js'\
  --word '../public/assets/frontend/js/page/*.js'\
  --word '../public/assets/frontend/js/page/*/*.js'\
  --include-selector ".modal-backdrop"\
  --include-selector ".modal-backdrop.show"\
  --include-selector ".modal-backdrop.fade" > ../public/assets/frontend/blogdefault/css/vendors/bootstrap.lite.min.css

cssdalek \
  --css '../public/assets/frontend/blogdefault/css/style.css'\
  --word '../public/assets/frontend/js/bootstrap/bootstrap.js'\
  --word '../template/frontend/blogdefault/*.html'\
  --word '../template/frontend/blogdefault/*/*.html'\
  --word '../template/frontend/blogdefault/*/*/*.html'\
  --word '../public/assets/frontend/js/*.js'\
  --word '../public/assets/frontend/js/page/*.js'\
  --word '../public/assets/frontend/js/page/*/*.js'\
  --include-selector ".progress-wrap.active-progress"\
  --include-id 'captchaImage' > ../public/assets/frontend/blogdefault/css/style.lite.min.css


# class="[^"]*-\{\{if