$(function(){
    $(".editProduct").click(
        function() {
            $("#productId").val($(this).attr('data-id'));
            $("#productName").val($(this).attr('data-name'));
            $("#productPrice").val($(this).attr('data-price'));
            $("#productAmount").val($(this).attr('data-amount'));
            $("#productImage").val($(this).attr('data-image'));
            $("#productCategory").val($(this).attr('data-category'));
        });
});

$(function(){
    $(".deleteProduct").click(
        function() {
            $(".productIdDel").val($(this).attr('data-id'));
            $("#productNameDel").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".editProductCharacteristics").click(
        function() {
            var id = $(this).attr('data-id');
            $("#allTables").children().hide();
            $("#product" + id).show();
            $("#btn" + id).show();
        });
});

$(function(){
    $(".createCharacteristic").click(
        function() {
            $('#productCharacteristics').modal('hide');
            $("#characteristicRelationCreate").val($(this).attr('data-relation'));
            if ($(this).attr('data-relation') !== '0') {
                $("#characteristicTypeCreate").hide();
                $("#characteristicTypeCreate").removeAttr("required");
            } else {
                $("#characteristicTypeCreate").show();
                $("#characteristicTypeCreate").attr("required", "");
            }
        });
});

$(function(){
    $(".editCharacteristic").click(
        function() {
            $("#characteristicId").val($(this).attr('data-id'));
            $("#characteristicName").val($(this).attr('data-name'));
            if ($(this).attr('data-relation') !== '0') {
                $("#characteristicType").hide();
                $("#characteristicType").removeAttr("required");
            } else {
                $("#characteristicType").show();
                $("#characteristicType").attr("required", "");
                $("#characteristicType").val($(this).attr('data-type'));
            }
            $("#characteristicRelationEdit").val($(this).attr('data-relation'));
        });
});

$(function(){
    $(".deleteCharacteristic").click(
        function() {
            $(".characteristicIdDel").val($(this).attr('data-id'));
            $("#characteristicNameDel").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".createCategory").click(
        function() {
            $("#categoryRelationCreate").val($(this).attr('data-relation'));
        });
});

$(function(){
    $(".editCategory").click(
        function() {
            $("#categoryId").val($(this).attr('data-id'));
            $("#categoryName").val($(this).attr('data-name'));
            $("#categoryRelationEdit").val($(this).attr('data-relation'));
        });
});

$(function(){
    $(".deleteCategory").click(
        function() {
            $(".categoryIdDel").val($(this).attr('data-id'));
            $("#categoryNameDel").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".createSet").click(
        function() {
            $('#productCharacteristics').modal('hide');
            $("#product").val($(this).attr('data-id'));
            $("#INTC").hide();
            $("#SETC").hide();
        });
});

$(function(){
    $("#setSelectorCreate").change(
        function() {
            var text = $(this).find(":selected").text();
            var selection = text.split(" (")[1].slice(0, -1);
            if (selection === "Список") {
                $("#INTC").hide();
                $("#INTC").removeAttr("required");
                $("#SETC").show();
                $("#SETC").children().hide();
                $("#SETC").children().removeAttr("required");
                $("#SETC" + $(this).find(":selected").val()).attr("required", "");
                $("#SETC" + $(this).find(":selected").val()).show();
            } else {
                $("#INTC").show();
                $("#INTC").attr("required", "");
                $("#SETC").hide();
                $("#SETC").children().removeAttr("required");
            }
        });
});

$(function(){
    $(".editSet").click(
        function() {
            $('#productCharacteristics').modal('hide');
            $("#setId").val($(this).attr('data-id'));
            if ($(this).attr('data-relation') !== "0") {
                $("#setSelectorEdit").val($(this).attr('data-relation'));
                $("#INTE").hide();
                $("#INTE").removeAttr("required");
                $("#SETE").show();
                $("#SETE").children().hide();
                $("#SETE").children().removeAttr("required");
                $("#SETE" + $(this).attr('data-relation')).attr("required", "");
                $("#SETE" + $(this).attr('data-relation')).show();
                $("#SETE" + $(this).attr('data-relation')).val($(this).attr('data-characteristic'));
            } else {
                $("#setSelectorEdit").val($(this).attr('data-characteristic'));
                $("#INTE").show();
                $("#INTE").attr("required", "");
                $("#ValueInt").val($(this).attr('data-value'));
                $("#SETE").hide();
                $("#SETE").children().removeAttr("required");
            }
        });
});

$(function(){
    $("#setSelectorEdit").change(
        function() {
            var text = $(this).find(":selected").text();
            var selection = text.split(" (")[1].slice(0, -1);
            if (selection === "Список") {
                $("#INTE").hide();
                $("#INTE").removeAttr("required");
                $("#SETE").show();
                $("#SETE").children().hide();
                $("#SETE").children().removeAttr("required");
                $("#SETE" + $(this).find(":selected").val()).attr("required", "");
                $("#SETE" + $(this).find(":selected").val()).show();
            } else {
                $("#INTE").show();
                $("#INTE").attr("required", "");
                $("#SETE").hide();
                $("#SETE").children().removeAttr("required");
            }
        });
});

$(function(){
    $(".deleteSet").click(
        function() {
            $('#productCharacteristics').modal('hide');
            $("#setIdDel").val($(this).attr('data-id'));
            $("#setCharacteristicDel").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".editStatus").click(
        function() {
            $("#statusId").val($(this).attr('data-id'));
            $("#statusStatus").val($(this).attr('data-status'));
        });
});

$(function(){
    $(".deleteStatus").click(
        function() {
            $(".statusIdDel").val($(this).attr('data-id'));
            $("#statusStatusDel").val($(this).attr('data-status'));
        });
});

$(function(){
    $(".editRole").click(
        function() {
            $("#roleId").val($(this).attr('data-id'));
            $("#roleName").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".deleteRole").click(
        function() {
            $(".roleIdDel").val($(this).attr('data-id'));
            $("#roleNameDel").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".editUser").click(
        function() {
            $("#userId").val($(this).attr('data-id'));
            $("#userEmail").val($(this).attr('data-email'));
            $("#userPhone").val($(this).attr('data-phone'));
            $("#userLastName").val($(this).attr('data-lastname'));
            $("#userName").val($(this).attr('data-name'));
            $("#userMiddleName").val($(this).attr('data-middlename'));
            $("#userGender").val($(this).attr('data-gender'));
            $("#userRole").val($(this).attr('data-role'));
        });
});

$(function(){
    $(".deleteUser").click(
        function() {
            $(".userIdDel").val($(this).attr('data-id'));
            $("#userEmailDel").val($(this).attr('data-email'));
        });
});

$(function(){
    $(".editOrder").click(
        function() {
            $("#orderId").val($(this).attr('data-id'));
            $("#orderDate").val($(this).attr('data-date'));
            $("#orderAddress").val($(this).attr('data-address'));
            $("#orderStatus").val($(this).attr('data-status'));
            $("#orderUser").val($(this).attr('data-user'));
        });
});

$(function(){
    $(".deleteOrder").click(
        function() {
            $("#orderIdDel").val($(this).attr('data-id'));
            $("#orderNumberDel").val(parseInt($(this).attr('data-id')) + 1000);
        });
});

$(function(){
    $(".editOrderPositions").click(
        function() {
            var id = $(this).attr('data-id');
            $("#allTablesPositions").children().hide();
            $("#position" + id).show();
            $("#btnp" + id).show();
        });
});

$(function(){
    $(".createPosition").click(
        function() {
            $('#orderPositions').modal('hide');
            $("#order").val($(this).attr('data-id'));
        });
});

$(function(){
    $("#positionSelectorCreate").change(
        function() {
            var text = $(this).find(":selected").text();
            var price = text.split(": ")[1].slice(0, -4);
            $("#positionPriceCreate").val(price);
            $("#positionAmountCreate").val(1);
        });
});

$(function(){
    $(".editPosition").click(
        function() {
            $('#orderPositions').modal('hide');
            $("#positionId").val($(this).attr('data-id'));
            $("#positionOrder").val($(this).attr('data-order'));
            $("#positionProduct").val($(this).attr('data-product'));
            $("#positionPrice").val($(this).attr('data-price'));
            $("#positionAmount").val($(this).attr('data-amount'));
        });
});

$(function(){
    $("#positionProduct").change(
        function() {
            var text = $(this).find(":selected").text();
            var price = text.split(": ")[1].slice(0, -4);
            $("#positionPrice").val(price);
            $("#positionAmount").val(1);
        });
});

$(function(){
    $(".deletePosition").click(
        function() {
            $('#orderPositions').modal('hide');
            $("#positionIdDel").val($(this).attr('data-id'));
            $("#positionProductDel").val($(this).attr('data-name'));
        });
});

$(function(){
    $(".editCard").click(
        function() {
            $("#cardId").val($(this).attr('data-id'));
            $("#cardNumber").val($(this).attr('data-number'));
            $("#cardDiscount").val($(this).attr('data-discount'));
            $("#cardDate").val($(this).attr('data-date'));
            $("#cardIdDel").val($(this).attr('data-id'));
        });
});

$(function(){
    $(".showImage").click(
        function() {
            $('#innerImage').attr('src', $(this).attr('data-src'));
        });
});

$(function(){
    $(".infoProduct").click(
        function() {
            var name = $(this).attr('data-name');
            $(".productName").text($(this).attr('data-name'));
            $(".productNumber").text("Номер: " + $(this).attr('data-number'));
            $("#productImage").attr("src", $(this).attr('data-image'));
            $(".productPrice").text("Цена: " + $(this).attr('data-price') + " руб");
            $(".productAmount").text("В наличии: " + $(this).attr('data-amount') + " шт");
            $(".productCategory").text($(this).attr('data-category'));
            $("#allCharacteristics").children().hide();
            $("#productCharacteristics" + $(this).attr('data-id')).show();
        });
});