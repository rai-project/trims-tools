
prettyNameMapping = <|
     "bvlc_alexnet_1.0"->"BVLC-AlexNet",
     "bvlc_googlenet_1.0"->"BVLC-GoogLeNet",
     "bvlc_reference_caffenet_1.0"->"BVLC-Reference-CaffeNet",
     "bvlc_reference_rcnn_ilsvrc13_1.0"->"BVLC-Reference-RCNN-ILSVRC13",
     "dpn68_1.0"->"DPN68",
     "dpn92_1.0"->"DPN92",
     "inception_3.0"->"Inception-v3",
     "inception_4.0"->"Inception-v4",
     "inceptionbn_21k_2.0"->"InceptionBN-21K-v2",
     "inception_bn_3.0"->"Inception-BN-v3",
     "inception_resnet_2.0"->"Inception-ResNet-v2",
     "locationnet_1.0"->"LocationNet",
     "network_in_network_1.0"->"NIN-v1",
     "o_resnet101_2.0"->"o-ResNet101-v2",
     "o_resnet152_2.0"->"ResNet152-v2",
     "o_vgg16_1.0"->"o-VGG16",
     "o_vgg19_1.0"->"o-VGG19",
     "resnet101_1.0"->"ResNet101-v1",
     "resnet101_2.0"->"ResNet101-v2",
     "resnet152_1.0"->"ResNet152-11k-v1",
     "resnet152_11k_1.0"->"ResNet152-11k",
     "resnet152_2.0"->"ResNet152-v2",
     "resnet18_2.0"->"ResNet18-v2",
     "resnet200_2.0"->"ResNet200-v2",
     "resnet269_2.0"->"ResNet269-v2",
     "resnet34_2.0"->"ResNet34-v2",
     "resnet50_1.0"->"ResNet50-v2",
     "resnet50_2.0"->"ResNet50-v2",
     "resnext101_1.0"->"ResNeXt101-v1",
     "resnext101_32x4d_1.0"->"ResNeXt101-32x4d",
     "resnext26_32x4d_1.0"->"ResNeXt26-32x4d",
     "resnext50_1.0"->"ResNeXt50",
     "resnext50_32x4d_1.0"->"ResNeXt50-32x4d",
     "squeezenet_1.0"->"SqueezeNet-v1.0",
     "squeezenet_1.1"->"SqueezeNet-v1.1",
     "vgg16_1.0"->"VGG16-v1",
     "vgg16_sod_1.0"->"VGG16_SOD",
     "vgg16_sos_1.0"->"VGG16_SOS",
     "vgg19_1.0"->"VGG19",
     "wrn50_2.0"->"WRN50-2",
     "xception_1.0"->"Xception"
|>;

prettyNameMapping = Join[prettyNameMapping, KeyValueMap[#2 -> #1 &, prettyNameMapping]];

prettyName[model_] := Lookup[prettyNameMapping, model];

