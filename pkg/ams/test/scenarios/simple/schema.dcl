schema {
    entity1 : {
        name: String,
        size: Number,
        public: Boolean,
        group: String
    },
    entity2 : {
        name: String,
        owners: String[],
        subentity: {
            subNumberField: Number
        }
    }
}