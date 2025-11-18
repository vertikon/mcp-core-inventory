/**
 * Schema Converter - Converts between JSON Schema, OpenAPI, and AsyncAPI formats
 * 
 * This module provides conversion utilities for different schema formats
 * used in the MCP-Hulk ecosystem.
 */

/**
 * Converts JSON Schema to OpenAPI Schema
 * @param {Object} jsonSchema - JSON Schema object
 * @returns {Object} OpenAPI Schema object
 */
function jsonSchemaToOpenAPI(jsonSchema) {
    if (!jsonSchema || typeof jsonSchema !== 'object') {
        throw new Error('Invalid JSON Schema provided');
    }

    const openAPI = {
        openapi: '3.0.0',
        info: {
            title: jsonSchema.title || 'API',
            version: jsonSchema.version || '1.0.0',
            description: jsonSchema.description || '',
        },
        paths: {},
        components: {
            schemas: {}
        }
    };

    // Convert JSON Schema properties to OpenAPI schema
    if (jsonSchema.properties) {
        const schemaName = jsonSchema.title || 'Schema';
        openAPI.components.schemas[schemaName] = convertSchemaProperties(jsonSchema);
    }

    return openAPI;
}

/**
 * Converts OpenAPI Schema to JSON Schema
 * @param {Object} openAPI - OpenAPI Schema object
 * @returns {Object} JSON Schema object
 */
function openAPIToJSONSchema(openAPI) {
    if (!openAPI || typeof openAPI !== 'object') {
        throw new Error('Invalid OpenAPI Schema provided');
    }

    const jsonSchema = {
        $schema: 'http://json-schema.org/draft-07/schema#',
        title: openAPI.info?.title || 'Schema',
        version: openAPI.info?.version || '1.0.0',
        description: openAPI.info?.description || '',
        type: 'object',
        properties: {}
    };

    // Convert OpenAPI components to JSON Schema
    if (openAPI.components?.schemas) {
        for (const [name, schema] of Object.entries(openAPI.components.schemas)) {
            jsonSchema.properties[name] = convertOpenAPISchemaToJSON(schema);
        }
    }

    return jsonSchema;
}

/**
 * Converts JSON Schema to AsyncAPI Schema
 * @param {Object} jsonSchema - JSON Schema object
 * @param {Object} options - Conversion options
 * @returns {Object} AsyncAPI Schema object
 */
function jsonSchemaToAsyncAPI(jsonSchema, options = {}) {
    if (!jsonSchema || typeof jsonSchema !== 'object') {
        throw new Error('Invalid JSON Schema provided');
    }

    const asyncAPI = {
        asyncapi: '2.6.0',
        info: {
            title: jsonSchema.title || 'Async API',
            version: jsonSchema.version || '1.0.0',
            description: jsonSchema.description || '',
        },
        channels: {},
        components: {
            schemas: {}
        }
    };

    // Convert JSON Schema to AsyncAPI message schema
    if (jsonSchema.properties) {
        const schemaName = jsonSchema.title || 'Message';
        asyncAPI.components.schemas[schemaName] = convertSchemaProperties(jsonSchema);
    }

    return asyncAPI;
}

/**
 * Converts AsyncAPI Schema to JSON Schema
 * @param {Object} asyncAPI - AsyncAPI Schema object
 * @returns {Object} JSON Schema object
 */
function asyncAPIToJSONSchema(asyncAPI) {
    if (!asyncAPI || typeof asyncAPI !== 'object') {
        throw new Error('Invalid AsyncAPI Schema provided');
    }

    const jsonSchema = {
        $schema: 'http://json-schema.org/draft-07/schema#',
        title: asyncAPI.info?.title || 'Schema',
        version: asyncAPI.info?.version || '1.0.0',
        description: asyncAPI.info?.description || '',
        type: 'object',
        properties: {}
    };

    // Convert AsyncAPI components to JSON Schema
    if (asyncAPI.components?.schemas) {
        for (const [name, schema] of Object.entries(asyncAPI.components.schemas)) {
            jsonSchema.properties[name] = convertAsyncAPISchemaToJSON(schema);
        }
    }

    return jsonSchema;
}

/**
 * Converts OpenAPI Schema to AsyncAPI Schema
 * @param {Object} openAPI - OpenAPI Schema object
 * @returns {Object} AsyncAPI Schema object
 */
function openAPIToAsyncAPI(openAPI) {
    const jsonSchema = openAPIToJSONSchema(openAPI);
    return jsonSchemaToAsyncAPI(jsonSchema);
}

/**
 * Converts AsyncAPI Schema to OpenAPI Schema
 * @param {Object} asyncAPI - AsyncAPI Schema object
 * @returns {Object} OpenAPI Schema object
 */
function asyncAPIToOpenAPI(asyncAPI) {
    const jsonSchema = asyncAPIToJSONSchema(asyncAPI);
    return jsonSchemaToOpenAPI(jsonSchema);
}

/**
 * Helper function to convert schema properties
 */
function convertSchemaProperties(schema) {
    const converted = {
        type: schema.type || 'object',
        properties: {},
        required: schema.required || []
    };

    if (schema.properties) {
        for (const [key, value] of Object.entries(schema.properties)) {
            converted.properties[key] = {
                type: value.type || 'string',
                description: value.description || '',
                example: value.example || null
            };
        }
    }

    return converted;
}

/**
 * Helper function to convert OpenAPI schema to JSON Schema
 */
function convertOpenAPISchemaToJSON(schema) {
    return {
        type: schema.type || 'object',
        properties: schema.properties || {},
        required: schema.required || []
    };
}

/**
 * Helper function to convert AsyncAPI schema to JSON Schema
 */
function convertAsyncAPISchemaToJSON(schema) {
    return {
        type: schema.type || 'object',
        properties: schema.properties || {},
        required: schema.required || []
    };
}

// Export functions
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        jsonSchemaToOpenAPI,
        openAPIToJSONSchema,
        jsonSchemaToAsyncAPI,
        asyncAPIToJSONSchema,
        openAPIToAsyncAPI,
        asyncAPIToOpenAPI
    };
}
