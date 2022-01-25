jQuery(function($){
    $("#product-img-c").click(function() {
      $("#overlay-product").fadeIn();　
    });
    $(".ol-close").click(function() {
      $("#overlay-product").fadeOut();
    });
    // バブリングを停止
    $(".overlay-inner").click(function(event){
      event.stopPropagation();
    });
  });