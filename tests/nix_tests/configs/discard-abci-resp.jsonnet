local config = import 'default.jsonnet';

config {
  'treasurenet_5005-1'+: {
    config+: {
      storage: {
        discard_abci_responses: true,
      },
    },
  },
}
