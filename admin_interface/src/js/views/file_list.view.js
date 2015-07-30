var $                  = require('jquery'),
    _                  = require('backbone'),
    Backbone           = require('backbone-crossdomain'),
    FilesCollection = require('../collections/files.collection')
    template           = require('../templates/file_list.hbs');

Backbone.$ = window.$

module.exports = Backbone.View.extend({
    el: $('#main'),
    template: template,

    initialize: function () {
        console.log("files initialize");
        var files = this;
        files.collection = new FilesCollection();
        files.collection.on("reset", this.render, this);
        files.collection.fetch({
            success: function () {
                files.render();
            }
        });
    },

    render: function () {
        this.$el.html(template({ files: this.collection.files() }));
        return this;
    }
});
