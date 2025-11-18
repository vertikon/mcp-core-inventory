/**
 * NATS Schema Generator - Generates NATS JetStream schemas and configurations
 * 
 * This module generates NATS subject schemas, streams, and JetStream configurations
 * for the MCP-Hulk messaging infrastructure.
 */

/**
 * Generates NATS subject schema configuration
 * @param {Object} config - Configuration object
 * @param {string} config.subject - NATS subject name
 * @param {Object} config.schema - JSON Schema for the subject
 * @param {Object} options - Generation options
 * @returns {Object} NATS subject configuration
 */
function generateNATSSubjectSchema(config, options = {}) {
    if (!config || !config.subject) {
        throw new Error('Subject name is required');
    }

    const subjectConfig = {
        subject: config.subject,
        schema: {
            type: 'json',
            schema: config.schema || {}
        },
        metadata: {
            name: config.subject,
            description: config.description || '',
            version: config.version || '1.0.0'
        }
    };

    // Add JetStream configuration if provided
    if (config.jetstream) {
        subjectConfig.jetstream = {
            stream: config.jetstream.stream || config.subject,
            consumer: config.jetstream.consumer || `${config.subject}-consumer`,
            ...config.jetstream
        };
    }

    return subjectConfig;
}

/**
 * Generates NATS JetStream stream configuration
 * @param {Object} config - Stream configuration
 * @param {string} config.name - Stream name
 * @param {string[]} config.subjects - List of subjects
 * @param {Object} options - Generation options
 * @returns {Object} JetStream stream configuration
 */
function generateNATSStream(config, options = {}) {
    if (!config || !config.name) {
        throw new Error('Stream name is required');
    }

    if (!config.subjects || !Array.isArray(config.subjects) || config.subjects.length === 0) {
        throw new Error('At least one subject is required');
    }

    const streamConfig = {
        name: config.name,
        subjects: config.subjects,
        retention: config.retention || 'limits',
        max_consumers: config.maxConsumers || -1,
        max_msgs: config.maxMsgs || -1,
        max_bytes: config.maxBytes || -1,
        max_age: config.maxAge || 0,
        storage: config.storage || 'file',
        num_replicas: config.numReplicas || 1,
        discard: config.discard || 'old',
        duplicate_window: config.duplicateWindow || 120000000000, // 2 minutes in nanoseconds
        compression: config.compression || 'none'
    };

    // Add placement if provided
    if (config.placement) {
        streamConfig.placement = config.placement;
    }

    // Add mirror if provided
    if (config.mirror) {
        streamConfig.mirror = config.mirror;
    }

    // Add sources if provided
    if (config.sources && Array.isArray(config.sources)) {
        streamConfig.sources = config.sources;
    }

    return streamConfig;
}

/**
 * Generates NATS JetStream consumer configuration
 * @param {Object} config - Consumer configuration
 * @param {string} config.stream - Stream name
 * @param {string} config.name - Consumer name
 * @param {Object} options - Generation options
 * @returns {Object} JetStream consumer configuration
 */
function generateNATSConsumer(config, options = {}) {
    if (!config || !config.stream) {
        throw new Error('Stream name is required');
    }

    if (!config.name) {
        throw new Error('Consumer name is required');
    }

    const consumerConfig = {
        stream_name: config.stream,
        name: config.name,
        durable_name: config.durableName || config.name,
        deliver_policy: config.deliverPolicy || 'all',
        ack_policy: config.ackPolicy || 'explicit',
        ack_wait: config.ackWait || 30000000000, // 30 seconds in nanoseconds
        max_deliver: config.maxDeliver || -1,
        filter_subject: config.filterSubject || '',
        replay_policy: config.replayPolicy || 'instant',
        rate_limit: config.rateLimit || 0,
        max_ack_pending: config.maxAckPending || 1000,
        max_waiting: config.maxWaiting || 512,
        headers_only: config.headersOnly || false
    };

    // Add flow control if provided
    if (config.flowControl !== undefined) {
        consumerConfig.flow_control = config.flowControl;
    }

    // Add heartbeat if provided
    if (config.heartbeat) {
        consumerConfig.heartbeat = config.heartbeat;
    }

    return consumerConfig;
}

/**
 * Generates complete NATS configuration from a schema definition
 * @param {Object} schemaDef - Schema definition
 * @param {Object} options - Generation options
 * @returns {Object} Complete NATS configuration
 */
function generateNATSConfig(schemaDef, options = {}) {
    if (!schemaDef || !schemaDef.subjects) {
        throw new Error('Schema definition with subjects is required');
    }

    const config = {
        jetstream: {
            enabled: options.enabled !== false,
            store_dir: options.storeDir || './data/jetstream',
            max_memory: options.maxMemory || '1GB',
            max_storage: options.maxStorage || '10GB'
        },
        streams: [],
        consumers: []
    };

    // Generate stream for each subject group
    const subjectGroups = groupSubjects(schemaDef.subjects);
    
    for (const [groupName, subjects] of Object.entries(subjectGroups)) {
        const streamConfig = generateNATSStream({
            name: groupName,
            subjects: subjects,
            retention: options.retention || 'limits',
            storage: options.storage || 'file'
        });

        config.streams.push(streamConfig);

        // Generate consumer for each stream
        const consumerConfig = generateNATSConsumer({
            stream: groupName,
            name: `${groupName}-consumer`,
            durableName: `${groupName}-consumer`
        });

        config.consumers.push(consumerConfig);
    }

    return config;
}

/**
 * Groups subjects by prefix/pattern
 */
function groupSubjects(subjects) {
    const groups = {};

    for (const subject of subjects) {
        const parts = subject.split('.');
        const groupName = parts.length > 1 ? parts.slice(0, -1).join('.') : 'default';
        
        if (!groups[groupName]) {
            groups[groupName] = [];
        }
        
        groups[groupName].push(subject);
    }

    return groups;
}

/**
 * Validates NATS subject name
 * @param {string} subject - Subject name to validate
 * @returns {boolean} True if valid
 */
function validateNATSSubject(subject) {
    if (!subject || typeof subject !== 'string') {
        return false;
    }

    // NATS subjects can contain: alphanumeric, dots, colons, underscores, dashes
    const pattern = /^[a-zA-Z0-9._:-]+$/;
    return pattern.test(subject);
}

// Export functions
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        generateNATSSubjectSchema,
        generateNATSStream,
        generateNATSConsumer,
        generateNATSConfig,
        validateNATSSubject
    };
}
