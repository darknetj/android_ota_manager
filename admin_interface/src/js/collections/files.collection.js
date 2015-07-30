var $        = require('jquery'),
    _               = require('backbone'),
    Backbone        = require('backbone-crossdomain'),
    File  = require("../models/file.model");

module.exports = Backbone.Collection.extend({

    model: File,

    urlRoot: 'http://localhost:8080',

    url: function() {
        return this.urlRoot + "/files.json";
    },

    files: function() {
        return this.toJSON()[0].result;
    }

});
