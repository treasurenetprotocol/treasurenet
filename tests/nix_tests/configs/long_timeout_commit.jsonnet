local default = import 'default.jsonnet';

default {
  'treasurenet_5005-1'+: {
    config+: {
      consensus+: {
        timeout_commit: '5s',
      },
    },
  },
}
