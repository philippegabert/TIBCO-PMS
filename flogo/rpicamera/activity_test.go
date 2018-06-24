{
  "name": "rpicamera",
  "title": "Capture image from Raspberry Pi Camera",
  "version": "0.0.1",
  "type": "flogo:activity",
  "description": "This activity allows you to capture image from Raspberry Pi camera.",
  "author": "Philippe GABERT <pgabert@tibco.com>",
  "ref": "github.com/philippegabert/flogo-contrib/activity/rpicamera",
  "homepage": "https://github.com/philippegabert/flogo-contrib/tree/rpicamera/activity/rpicamera",
  "inputs":[
    {
      "name": "deviceID",
      "type": "integer",
	  "required":"true"
    },
	{
      "name": "picWidth",
      "type": "integer"
    },
	{
      "name": "picHeight",
      "type": "integer"
    },
	{
      "name": "outputType",
      "type": "string",
	  "allowed" : ["file","base64"],
	  "value": "file"
    },
	{
      "name": "folderOut",
      "type": "string"
    }
  ],
  "outputs": [
   	{
      "name": "picFile",
      "type": "string"
    },
	{
      "name": "base64",
      "type": "string"
    }
  ]
}